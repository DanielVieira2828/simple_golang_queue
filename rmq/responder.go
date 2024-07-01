package rmq

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Responder struct {
	channels sync.Map
}

func (r *Responder) CreateResponseChannel() string {
	requestID := uuid.NewString()

	channel := make(chan ConsumerResponse)

	r.channels.Store(requestID, channel)

	return requestID
}

func (r *Responder) DeleteResponseChannel(requestID string) {
	r.channels.Delete(requestID)
}

func (r *Responder) WaitForResponse(requestID string, timeoutTime time.Duration) (ConsumerResponse, error) {
	channelInterface, ok := r.channels.Load(requestID)
	if !ok {
		return ConsumerResponse{}, fmt.Errorf("nenhum cliente está esperando resposta para a requisição %s", requestID)
	}

	defer r.channels.Delete(requestID)

	channel, ok := channelInterface.(chan ConsumerResponse)
	if !ok {
		return ConsumerResponse{}, fmt.Errorf("a requisição com o ID %s possui um tipo inválido", requestID)
	}

	defer close(channel)

	var response ConsumerResponse
	select {
	case response = <-channel:

		return response, nil

	case <-time.After(timeoutTime):

		return ConsumerResponse{
			Error: "a requisição atingiu o tempo limite",
		}, errors.New("a requisição atingiu o tempo limite")
	}
}

func (r *Responder) SendResponse(requestID string, response ConsumerResponse) {
	channelInterface, ok := r.channels.Load(requestID)
	if !ok {
		log.Printf("nenhum cliente está esperando resposta para a requisição %s\n", requestID)
		return
	}

	channel, ok := channelInterface.(chan ConsumerResponse)
	if !ok {
		log.Printf("a requisição com o ID %s possui um tipo inválido\n", requestID)
		return
	}

	channel <- response
}

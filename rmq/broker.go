package rmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adjust/rmq/v5"
)

type Broker struct {
	Queue     rmq.Queue
	Responder *Responder
}

func (b *Broker) MakeRequest() (ConsumerResponse, error) {

	requestID := b.Responder.CreateResponseChannel()

	request := ConsumerRequest{
		RequestId: requestID,
	}

	requestJSONBytes, err := json.Marshal(request)
	if err != nil {
		b.Responder.DeleteResponseChannel(requestID)
		return ConsumerResponse{}, fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	err = b.Queue.PublishBytes(requestJSONBytes)
	if err != nil {
		b.Responder.DeleteResponseChannel(requestID)
		return ConsumerResponse{}, fmt.Errorf("erro ao publicar requisição na fila: %w", err)
	}

	return b.Responder.WaitForResponse(requestID, time.Second*30)
}

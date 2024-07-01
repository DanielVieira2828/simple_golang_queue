package rmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/google/uuid"
)

type Consumer struct {
	responder *Responder
}

type ConsumerRequest struct {
	RequestId string
}

type ConsumerResponse struct {
	Uuid  string `json:"uuid"`
	Error string `json:"error,omitempty"`
}

func (c Consumer) Consume(delivery rmq.Delivery) {

	var request ConsumerRequest
	err := json.Unmarshal([]byte(delivery.Payload()), &request)
	if err != nil {
		log.Printf("erro ao desserializar requisição: %s\n", err.Error())
		delivery.Reject()
		return
	}

	time.Sleep(2 * time.Second)

	response := ConsumerResponse{
		Uuid: uuid.NewString(),
	}
	log.Printf("resposta gerada: %s\n", response.Uuid)

	c.responder.SendResponse(request.RequestId, response)

	if err := delivery.Ack(); err != nil {
		log.Printf("erro ao confirmar recebimento da mensagem: %s\n", err.Error())
	}
}

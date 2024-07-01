package rmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adjust/rmq/v5"
)

const (
	RedisDatabase = 1

	PrefetchLimit = 4

	PollDuration = time.Second
)

func Run() (conn rmq.Connection, responder *Responder, queue rmq.Queue, err error) {

	errorChannel := make(chan error)

	go func() {
		for err := range errorChannel {
			log.Printf("erro do RabbitMQ: %s\n", err.Error())
		}
	}()

	conn, err = rmq.OpenConnection("simple_queue_service", "tcp", "localhost:6379", RedisDatabase, errorChannel)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao conectar com o RabbitMQ: %w", err)
	}

	queue, err = conn.OpenQueue("default_queue")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao abrir a fila: %w", err)
	}

	if err := queue.StartConsuming(PrefetchLimit, PollDuration); err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao iniciar o consumo de mensagens: %w", err)
	}

	responder = &Responder{
		channels: sync.Map{},
	}

	if _, err := queue.AddConsumer("default_consumer", Consumer{responder: responder}); err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao adicionar consumidor à fila: %w", err)
	}

	log.Println("fila iniciada...")

	return conn, responder, queue, nil
}

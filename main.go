package main

import (
	"log"

	"github.com/savi2w/simple-queue/rmq"
	"github.com/savi2w/simple-queue/server"
)

func main() {

	conn, responder, queue, err := rmq.Run()
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() {
		channel := conn.StopAllConsuming()
		<-channel
	}()

	broker := &rmq.Broker{
		Queue:     queue,
		Responder: responder,
	}

	if err := server.Run(broker); err != nil {
		log.Println(err.Error())
		return
	}
}

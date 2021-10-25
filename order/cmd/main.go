package main

import (
	"fmt"
	"os"

	"github.com/antsla/order/pkg/broker"
	"github.com/antsla/order/pkg/datastore"
	"github.com/antsla/order/transport"
	"github.com/rs/zerolog/log"
)

func main() {
	db := datastore.InitDB()
	fmt.Println(1)
	producer := broker.InitKafkaProducer()
	fmt.Println(2)

	server := transport.NewServer(db, producer)
	fmt.Println("server is starting...")
	err := server.Start()
	if err != nil {
		log.Error().Err(err).Msg("Server hasn't been started.")
		os.Exit(1)
	}
}

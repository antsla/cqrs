package main

import (
	"fmt"
	"os"

	"github.com/antsla/goods/pkg/broker"
	"github.com/antsla/goods/pkg/datastore"
	"github.com/antsla/goods/transport"
	"github.com/rs/zerolog/log"
)

func main() {
	db := datastore.InitDB()
	producer := broker.InitKafkaProducer()

	server := transport.NewServer(db, producer)
	fmt.Println("server is starting...")
	err := server.Start()
	if err != nil {
		log.Error().Err(err).Msg("Server hasn't been started.")
		os.Exit(1)
	}
}

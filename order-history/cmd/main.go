package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/antsla/order-history/pkg/broker"
	"github.com/antsla/order-history/pkg/datastore"
	"github.com/antsla/order-history/transport"
	"github.com/rs/zerolog/log"
)

func main() {
	db := datastore.InitDB()

	handlers := map[string]sarama.ConsumerGroupHandler{
		os.Getenv("GOODS_CREATED_TOPIC"): broker.BuildGoodsHandler(db),
		os.Getenv("ORDER_CREATED_TOPIC"): broker.BuildOrderHandler(db),
	}
	broker.RunConsumers(context.Background(), handlers)

	server := transport.NewServer(db)
	fmt.Println("server is starting...")
	err := server.Start()
	if err != nil {
		log.Error().Err(err).Msg("Server hasn't been started.")
		os.Exit(1)
	}
}

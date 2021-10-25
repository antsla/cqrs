package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type GoodsCreatedEvent struct {
	Data struct {
		ID        int64     `json:"id"`
		OrderID   int64     `json:"order_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type GoodsHandler struct {
	db *pgxpool.Pool
}

func BuildGoodsHandler(db *pgxpool.Pool) GoodsHandler {
	return GoodsHandler{db: db}
}

func (gh GoodsHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		gce := GoodsCreatedEvent{}
		err := json.Unmarshal(msg.Value, &gce)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been handled.")
			session.MarkMessage(msg, "")
			continue
		}

		_, err = gh.db.Exec(context.Background(), `INSERT INTO goods (id, order_id, created_at) VALUES ($1, $2, $3)`, gce.Data.ID, gce.Data.OrderID, gce.Data.CreatedAt)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been inserted.")
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (GoodsHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (GoodsHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

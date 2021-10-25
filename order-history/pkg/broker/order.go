package broker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type OrderCreatedEvent struct {
	Data struct {
		ID        int64     `json:"id"`
		UserID    int64     `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

type OrderHandler struct {
	db *pgxpool.Pool
}

func BuildOrderHandler(db *pgxpool.Pool) OrderHandler {
	return OrderHandler{db: db}
}

func (oh OrderHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		oce := OrderCreatedEvent{}
		err := json.Unmarshal(msg.Value, &oce)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been handled.")
			session.MarkMessage(msg, "")
			continue
		}

		_, err = oh.db.Exec(context.Background(), `INSERT INTO orders (id, user_id, created_at) VALUES ($1, $2, $3)`, oce.Data.ID, oce.Data.UserID, oce.Data.CreatedAt)
		if err != nil {
			log.Error().Err(err).Msg("Event hasn't been inserted.")
		}

		session.MarkMessage(msg, "")
	}

	return nil
}

func (OrderHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (OrderHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

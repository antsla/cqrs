package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/antsla/goods/pkg/model"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router        *mux.Router
	db            *pgxpool.Pool
	kafkaProducer sarama.SyncProducer
}

func NewServer(db *pgxpool.Pool, kafkaProducer sarama.SyncProducer) Server {
	s := Server{}
	s.kafkaProducer = kafkaProducer
	s.db = db
	s.router = mux.NewRouter()

	s.router.HandleFunc("/v1/goods", s.CreateGoodsV1).Methods(http.MethodPost)

	return s
}

func (s Server) Start() error {
	return http.ListenAndServe(":"+os.Getenv("HTTP_BIND"), s.router)
}

func (s Server) CreateGoodsV1(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("Data hasn't been parsed.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderID := r.Form.Get("orderId")

	goods := model.Goods{}
	err = s.db.QueryRow(context.Background(), `INSERT INTO goods (order_id, created_at) VALUES ($1, NOW()) RETURNING id, order_id, created_at`, orderID).Scan(&goods.ID, &goods.OrderID, &goods.CreatedAt)
	if err != nil {
		log.Error().Err(err).Msg("Goods hasn't been created.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg := model.CreatedGoodsMsg{Data: goods}
	msgStr, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("Message hasn't been marshaled.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	producerMsg := &sarama.ProducerMessage{Topic: os.Getenv("GOODS_CREATED_TOPIC"), Value: sarama.StringEncoder(msgStr)}
	_, _, err = s.kafkaProducer.SendMessage(producerMsg)
	if err != nil {
		log.Error().Err(err).Msg("Message hasn't been sent.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

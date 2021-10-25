package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/antsla/order-history/pkg/model"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router *mux.Router
	db     *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) Server {
	s := Server{}
	s.db = db
	s.router = mux.NewRouter()

	s.router.HandleFunc("/v1/order-history", s.GetOrderHistoryV1).Methods(http.MethodGet)

	return s
}

func (s Server) Start() error {
	return http.ListenAndServe(":"+os.Getenv("HTTP_BIND"), s.router)
}

func (s Server) GetOrderHistoryV1(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("Data hasn't been parsed.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	threshold := r.Form.Get("threshold")
	offset := r.Form.Get("offset")
	limit := r.Form.Get("limit")

	rows, err := s.db.Query(context.Background(), `SELECT orders.id, orders.user_id, orders.created_at FROM orders
    	INNER JOIN goods ON goods.order_id = orders.id 
		GROUP BY orders.id 
		HAVING COUNT(goods.id) > $1
		LIMIT $2 OFFSET $3`, threshold, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Goods haven't been got.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	data := make([]model.Order, 0)
	for rows.Next() {
		o := model.Order{}
		sErr := rows.Scan(&o.ID, &o.UserID, &o.CreatedAt)
		if sErr != nil {
			log.Error().Err(err).Msg("Reading error.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		data = append(data, o)
	}

	ordersRsp := model.OrdersResponse{Data: data}
	response, err := json.Marshal(ordersRsp)
	if err != nil {
		log.Error().Err(err).Msg("Response hasn't been marshaled.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

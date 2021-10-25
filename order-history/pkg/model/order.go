package model

import "time"

type OrdersResponse struct {
	Data []Order `json:"data"`
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

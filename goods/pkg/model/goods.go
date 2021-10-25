package model

import "time"

type Goods struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatedGoodsMsg struct {
	Data Goods `json:"data"`
}

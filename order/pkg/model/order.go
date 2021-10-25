package model

import "time"

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatedOrderMsg struct {
	Data Order `json:"data"`
}

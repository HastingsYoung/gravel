package models

import (
	"time"
)

type Deal struct {
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
	Total     float64 `json:"total"`
	Timestamp int64   `json:"timestamp"`
}

func NewDeal(price, amount float64) *Deal {
	return &Deal{
		Price:     price,
		Amount:    amount,
		Total:     price * amount,
		Timestamp: time.Now().Unix(),
	}
}

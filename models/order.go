package models

import (
	"github.com/satori/go.uuid"
	"time"
)

const (
	ORDER_TYPE_ASK = "ASK"
	ORDER_TYPE_BID = "BID"
)

type Order struct {
	OrderId   string  `json:"order_id"`
	Market    string  `json:"market"`
	Type      string  `json:"type"`
	Stock     *Link   `json:"stock"`
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
	Total     float64 `json:"total"`
	Timestamp int64   `json:"timestamp"`
	Index     int     `json:"-"`
}

func NewOrder(market, tp string, stock *Stock, price, amount float64) *Order {
	return &Order{
		OrderId:   uuid.NewV4().String(),
		Market:    market,
		Type:      tp,
		Stock:     stock.RefLink(),
		Price:     price,
		Amount:    amount,
		Total:     price * amount,
		Timestamp: time.Now().Unix(),
		Index:     -1, // initialise index to -1 for safety
	}
}
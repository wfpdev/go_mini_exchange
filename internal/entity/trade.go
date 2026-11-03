package entity

import "time"

type Trade struct {
	ID          int64     `json:"id"`
	StockCode   string    `json:"stock_code"`
	BuyOrderID  int64     `json:"buy_order_id"`
	SellOrderID int64     `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	ExecutedAt  time.Time `json:"executed_at"`
}

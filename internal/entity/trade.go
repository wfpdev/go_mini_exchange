package entity

import "time"

type Trade struct {
	ID          int64
	StockCode   string
	BuyOrderID  int64
	SellOrderID int64
	Price       float64
	Quantity    int
	ExecutedAt  time.Time
}

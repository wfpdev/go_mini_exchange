package entity

import "time"

type OrderSide string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"
)

type OrderStatus string

const (
	Pending         OrderStatus = "PENDING"
	Open            OrderStatus = "OPEN"
	PartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	Filled          OrderStatus = "FILLED"
)

type Order struct {
	ID        int64
	UserID    int64
	StockCode string
	Side      OrderSide
	Price     float64

	OriginalQuantity  int
	FilledQuantity    int
	RemainingQuantity int

	Status OrderStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

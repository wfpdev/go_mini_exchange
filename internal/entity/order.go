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
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	StockCode string    `json:"stock_code"`
	Side      OrderSide `json:"side"`
	Price     float64   `json:"price"`

	OriginalQuantity  int `json:"original_quantity"`
	FilledQuantity    int `json:"filled_quantity"`
	RemainingQuantity int `json:"remaining_quantity"`

	Status OrderStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserID    int64   `json:"user_id"`
	StockCode string  `json:"stock_code"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

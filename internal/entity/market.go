package entity

import "time"

type Ticker struct {
	StockCode     string  `json:"stock_code"`
	Name          string  `json:"name"`
	LastPrice     float64 `json:"last_price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Volume        int     `json:"volume"`
}

type OrderBookLevel struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderBookSnapshot struct {
	StockCode string           `json:"stock_code"`
	Buy       []OrderBookLevel `json:"buy"`
	Sell      []OrderBookLevel `json:"sell"`
	UpdatedAt time.Time        `json:"updated_at"`
}

package entity

type Stock struct {
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	InitialPrice    float64 `json:"initial_price"`
	LastTradedPrice float64 `json:"last_traded_price"`
}

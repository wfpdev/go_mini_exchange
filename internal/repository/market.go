package repository

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
)

type MarketRepository interface {
	GetTicker(stockCode string) (*entity.Ticker, error)
	GetOrderBook(stockCode string) (*entity.OrderBookSnapshot, error)
}

type InMemoryMarketRepository struct {
	mu       sync.Mutex
	exchange *entity.Exchange
}

func NewInMemoryMarketRepository(exchange *entity.Exchange) *InMemoryMarketRepository {
	return &InMemoryMarketRepository{exchange: exchange}
}

func (r *InMemoryMarketRepository) GetTicker(stockCode string) (*entity.Ticker, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stockCode = strings.ToUpper(strings.TrimSpace(stockCode))
	stock, exists := r.exchange.Stocks[stockCode]
	if !exists {
		return nil, errors.New("stock not found")
	}

	volume := 0
	for _, trade := range r.exchange.Trades {
		if trade.StockCode == stockCode {
			volume += trade.Quantity
		}
	}

	change := stock.LastTradedPrice - stock.InitialPrice
	changePercent := 0.0
	if stock.InitialPrice > 0 {
		changePercent = (change / stock.InitialPrice) * 100
	}

	return &entity.Ticker{
		StockCode:     stock.Code,
		Name:          stock.Name,
		LastPrice:     stock.LastTradedPrice,
		Change:        change,
		ChangePercent: changePercent,
		Volume:        volume,
	}, nil
}

func (r *InMemoryMarketRepository) GetOrderBook(stockCode string) (*entity.OrderBookSnapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stockCode = strings.ToUpper(strings.TrimSpace(stockCode))
	if _, exists := r.exchange.Stocks[stockCode]; !exists {
		return nil, errors.New("stock not found")
	}

	book := r.exchange.OrderBooks[stockCode]
	buy := aggregateLevels(book.BuyOrders)
	sell := aggregateLevels(book.SellOrders)

	sort.Slice(buy, func(i, j int) bool { return buy[i].Price > buy[j].Price })
	sort.Slice(sell, func(i, j int) bool { return sell[i].Price < sell[j].Price })

	return &entity.OrderBookSnapshot{
		StockCode: stockCode,
		Buy:       buy,
		Sell:      sell,
		UpdatedAt: time.Now(),
	}, nil
}

func aggregateLevels(orders []*entity.Order) []entity.OrderBookLevel {
	quantities := make(map[float64]int)
	for _, order := range orders {
		if order.RemainingQuantity > 0 {
			quantities[order.Price] += order.RemainingQuantity
		}
	}

	levels := make([]entity.OrderBookLevel, 0, len(quantities))
	for price, quantity := range quantities {
		levels = append(levels, entity.OrderBookLevel{Price: price, Quantity: quantity})
	}
	return levels
}

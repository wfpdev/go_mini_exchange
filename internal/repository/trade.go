package repository

import (
	"strings"
	"sync"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
)

type TradeRepository interface {
	GetTrades(stockCode string) []*entity.Trade
}

type InMemoryTradeRepository struct {
	mu       sync.Mutex
	exchange *entity.Exchange
}

func NewInMemoryTradeRepository(exchange *entity.Exchange) *InMemoryTradeRepository {
	return &InMemoryTradeRepository{exchange: exchange}
}

func (r *InMemoryTradeRepository) GetTrades(stockCode string) []*entity.Trade {
	r.mu.Lock()
	defer r.mu.Unlock()

	stockCode = strings.ToUpper(strings.TrimSpace(stockCode))
	trades := make([]*entity.Trade, 0)

	for _, trade := range r.exchange.Trades {
		if stockCode != "" && !strings.EqualFold(trade.StockCode, stockCode) {
			continue
		}
		trades = append(trades, trade)
	}

	return trades
}

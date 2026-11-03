package usecase

import (
	"strings"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
)

type TradeUseCase interface {
	GetTrades(stockCode string) []*entity.Trade
}

type tradeUseCase struct {
	repository repository.TradeRepository
}

func NewTradeUseCase(repository repository.TradeRepository) TradeUseCase {
	return &tradeUseCase{repository: repository}
}

func (uc *tradeUseCase) GetTrades(stockCode string) []*entity.Trade {
	return uc.repository.GetTrades(strings.ToUpper(strings.TrimSpace(stockCode)))
}

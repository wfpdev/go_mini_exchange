package usecase

import (
	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
)

type MarketUseCase interface {
	GetTicker(stockCode string) (*entity.Ticker, error)
	GetOrderBook(stockCode string) (*entity.OrderBookSnapshot, error)
}

type marketUseCase struct {
	repository repository.MarketRepository
}

func NewMarketUseCase(repository repository.MarketRepository) MarketUseCase {
	return &marketUseCase{repository: repository}
}

func (uc *marketUseCase) GetTicker(stockCode string) (*entity.Ticker, error) {
	return uc.repository.GetTicker(stockCode)
}

func (uc *marketUseCase) GetOrderBook(stockCode string) (*entity.OrderBookSnapshot, error) {
	return uc.repository.GetOrderBook(stockCode)
}

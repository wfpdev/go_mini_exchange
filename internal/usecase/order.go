package usecase

import (
	"strings"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
)

type OrderUseCase interface {
	CreateOrder(request entity.CreateOrderRequest) (*entity.Order, error)
	GetOrders(userID *int64, stockCode string, status string) []*entity.Order
}

type orderUseCase struct {
	repository repository.OrderRepository
}

func NewCreateOrderUseCase(repository repository.OrderRepository) OrderUseCase {
	return &orderUseCase{repository: repository}
}

func (uc *orderUseCase) CreateOrder(request entity.CreateOrderRequest) (*entity.Order, error) {
	return uc.repository.CreateOrder(
		request.UserID,
		strings.ToUpper(strings.TrimSpace(request.StockCode)),
		entity.OrderSide(strings.ToUpper(strings.TrimSpace(request.Side))),
		request.Price,
		request.Quantity,
	)
}

func (uc *orderUseCase) GetOrders(userID *int64, stockCode string, status string) []*entity.Order {
	status = strings.ToUpper(strings.TrimSpace(status))
	status = strings.ReplaceAll(status, "-", "_")
	status = strings.ReplaceAll(status, " ", "_")

	return uc.repository.GetOrders(
		userID,
		strings.ToUpper(strings.TrimSpace(stockCode)),
		entity.OrderStatus(status),
	)
}

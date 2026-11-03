package repository

import (
	"strings"
	"sync"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
)

type OrderRepository interface {
	CreateOrder(userID int64, stockCode string, side entity.OrderSide, price float64, quantity int) (*entity.Order, error)
	GetOrders(userID *int64, stockCode string, status entity.OrderStatus) []*entity.Order
}

type InMemoryOrderRepository struct {
	mu       sync.Mutex
	exchange *entity.Exchange
	orders   []*entity.Order
}

func NewInMemoryOrderRepository(exchange *entity.Exchange) *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		exchange: exchange,
		orders:   make([]*entity.Order, 0),
	}
}

func (r *InMemoryOrderRepository) CreateOrder(userID int64, stockCode string, side entity.OrderSide, price float64, quantity int) (*entity.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, err := r.exchange.CreateOrder(userID, stockCode, side, price, quantity)
	if err != nil {
		return nil, err
	}

	r.orders = append(r.orders, order)
	return order, nil
}

func (r *InMemoryOrderRepository) GetOrders(userID *int64, stockCode string, status entity.OrderStatus) []*entity.Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	stockCode = strings.ToUpper(strings.TrimSpace(stockCode))
	orders := make([]*entity.Order, 0)

	for _, order := range r.orders {
		if userID != nil && order.UserID != *userID {
			continue
		}
		if stockCode != "" && !strings.EqualFold(order.StockCode, stockCode) {
			continue
		}
		if status != "" && order.Status != status {
			continue
		}

		orders = append(orders, order)
	}

	return orders
}

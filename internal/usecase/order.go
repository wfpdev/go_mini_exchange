package usecase

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
)

type OrderUseCase interface {
	CreateOrder(request entity.CreateOrderRequest) (*entity.Order, error)
	GetOrders(userID *int64, stockCode string, status string) []*entity.Order
}

type orderUseCase struct {
	repository repository.OrderRepository
	cache      *cache.Cache
	mutex      sync.Mutex
}

func NewCreateOrderUseCase(repository repository.OrderRepository) OrderUseCase {
	return &orderUseCase{
		repository: repository,
		cache:      cache.New(1*time.Minute, 5*time.Minute),
	}
}

func (uc *orderUseCase) CreateOrder(request entity.CreateOrderRequest) (*entity.Order, error) {
	key, err := createOrderIdempotencyKey(request)
	if err != nil {
		return nil, err
	}

	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if cached, found := uc.cache.Get(key); found {
		if order, ok := cached.(*entity.Order); ok {
			return order, nil
		}
	}

	order, err := uc.repository.CreateOrder(
		request.UserID,
		strings.ToUpper(strings.TrimSpace(request.StockCode)),
		entity.OrderSide(strings.ToUpper(strings.TrimSpace(request.Side))),
		request.Price,
		request.Quantity,
	)
	if err != nil {
		return nil, err
	}

	uc.cache.Set(key, order, cache.DefaultExpiration)
	return order, nil
}

func createOrderIdempotencyKey(request entity.CreateOrderRequest) (string, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(body)
	return fmt.Sprintf("create-order:%x", sum), nil
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

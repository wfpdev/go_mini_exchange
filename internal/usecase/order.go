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
	"github.com/wfpdev/go_mini_exchange/internal/messaging"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
)

type OrderUseCase interface {
	CreateOrder(request entity.CreateOrderRequest) (*entity.Order, error)
	ProcessOrder(orderID int64) (*entity.Order, error)
	GetOrders(userID *int64, stockCode string, status string) []*entity.Order
}

type orderUseCase struct {
	repository repository.OrderRepository
	cache      *cache.Cache
	publisher  messaging.KafkaOrderPublisher
	mutex      sync.Mutex
}

func NewCreateOrderUseCase(repository repository.OrderRepository, publishers ...messaging.KafkaOrderPublisher) OrderUseCase {
	var publisher messaging.KafkaOrderPublisher
	if len(publishers) > 0 {
		publisher = publishers[0]
	}
	return &orderUseCase{
		repository: repository,
		cache:      cache.New(1*time.Minute, 5*time.Minute),
		publisher:  publisher,
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
			cachedOrder := *order
			return &cachedOrder, nil
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

	pendingOrder := *order

	if uc.publisher != nil {
		if err := uc.publisher.Publish(order.ID); err != nil {
			return nil, fmt.Errorf("publish create order event: %w", err)
		}
	}

	uc.cache.Set(key, &pendingOrder, cache.DefaultExpiration)
	return &pendingOrder, nil
}

func createOrderIdempotencyKey(request entity.CreateOrderRequest) (string, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(body)
	return fmt.Sprintf("create-order:%x", sum), nil
}

func (uc *orderUseCase) ProcessOrder(orderID int64) (*entity.Order, error) {
	return uc.repository.ProcessOrder(orderID)
}

func (uc *orderUseCase) GetOrders(userID *int64, stockCode string, status string) []*entity.Order {
	status = strings.ToUpper(strings.TrimSpace(status))
	return uc.repository.GetOrders(
		userID,
		strings.ToUpper(strings.TrimSpace(stockCode)),
		entity.OrderStatus(status),
	)
}

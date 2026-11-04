package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/usecase"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUseCase
}

func NewOrderHandler(orderUsecase usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var request entity.CreateOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	order, err := h.orderUsecase.CreateOrder(request)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		return
	}
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	var userID *int64
	if value := r.URL.Query().Get("user_id"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 0 {
			writeError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		userID = &parsed
	}

	stockCode := r.URL.Query().Get("stock_code")
	status := r.URL.Query().Get("status")
	orders := h.orderUsecase.GetOrders(userID, stockCode, status)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(orders)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

type OrderConsumerHandler struct {
	orderUseCase usecase.OrderUseCase
}

func NewOrderConsumerHandler(orderUseCase usecase.OrderUseCase) *OrderConsumerHandler {
	return &OrderConsumerHandler{orderUseCase: orderUseCase}
}

func (h *OrderConsumerHandler) HandleCreateOrder(orderID int64) error {
	_, err := h.orderUseCase.ProcessPendingOrder(orderID)
	return err
}

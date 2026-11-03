package handler

import (
	"encoding/json"
	"net/http"

	"github.com/wfpdev/go_mini_exchange/internal/usecase"
)

type TradeHandler struct {
	tradeUseCase usecase.TradeUseCase
}

func NewTradeHandler(tradeUseCase usecase.TradeUseCase) *TradeHandler {
	return &TradeHandler{tradeUseCase: tradeUseCase}
}

func (h *TradeHandler) GetTrades(w http.ResponseWriter, r *http.Request) {
	stockCode := r.URL.Query().Get("stock_code")
	trades := h.tradeUseCase.GetTrades(stockCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(trades)
}

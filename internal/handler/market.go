package handler

import (
	"encoding/json"
	"net/http"

	"github.com/wfpdev/go_mini_exchange/internal/usecase"
)

type MarketHandler struct {
	marketUseCase usecase.MarketUseCase
}

func NewMarketHandler(marketUseCase usecase.MarketUseCase) *MarketHandler {
	return &MarketHandler{marketUseCase: marketUseCase}
}

func (h *MarketHandler) GetTicker(w http.ResponseWriter, r *http.Request) {
	stockCode := marketStockCode(r)
	if stockCode == "" {
		writeError(w, http.StatusBadRequest, "stock_code is required")
		return
	}

	ticker, err := h.marketUseCase.GetTicker(stockCode)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ticker)
}

func (h *MarketHandler) GetOrderBook(w http.ResponseWriter, r *http.Request) {
	stockCode := marketStockCode(r)
	if stockCode == "" {
		writeError(w, http.StatusBadRequest, "stock_code is required")
		return
	}

	snapshot, err := h.marketUseCase.GetOrderBook(stockCode)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(snapshot)
}

func marketStockCode(r *http.Request) string {
	stockCode := r.URL.Query().Get("stock_code")
	if stockCode == "" {
		stockCode = r.URL.Query().Get("stockCode")
	}
	return stockCode
}

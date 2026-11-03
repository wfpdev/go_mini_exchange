package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/usecase"
)

type MarketHandler struct {
	marketUseCase usecase.MarketUseCase
	tradeUseCase  usecase.TradeUseCase
}

func NewMarketHandler(marketUseCase usecase.MarketUseCase, tradeUseCase usecase.TradeUseCase) *MarketHandler {
	return &MarketHandler{marketUseCase: marketUseCase, tradeUseCase: tradeUseCase}
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
	return strings.ToUpper(strings.TrimSpace(stockCode))
}

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *MarketHandler) StreamMarket(w http.ResponseWriter, r *http.Request) {
	stockCode := marketStockCode(r)
	if stockCode == "" {
		writeError(w, http.StatusBadRequest, "stock_code is required")
		return
	}
	if _, err := h.marketUseCase.GetTicker(stockCode); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	topicParam := strings.TrimSpace(r.URL.Query().Get("topics"))

	tickerEnabled := true
	tradeEnabled := true
	if topicParam != "" {
		tickerEnabled = false
		tradeEnabled = false
		for _, topic := range strings.Split(topicParam, ",") {
			switch strings.TrimSpace(topic) {
			case "market.ticker":
				tickerEnabled = true
			case "market.trade":
				tradeEnabled = true
			}
		}
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	seenTradeIDs := make(map[int64]struct{})
	for _, trade := range h.tradeUseCase.GetTrades(stockCode) {
		seenTradeIDs[trade.ID] = struct{}{}
	}

	var lastTicker *entity.Ticker
	ticker, err := h.marketUseCase.GetTicker(stockCode)
	if err == nil && tickerEnabled {
		_ = conn.WriteJSON(map[string]interface{}{"topic": "market.ticker", "data": ticker})
		lastTicker = ticker
	}

	tickerInterval := time.NewTicker(time.Second)
	defer tickerInterval.Stop()

	for {
		select {
		case <-done:
			return
		case <-tickerInterval.C:
			if tickerEnabled {
				current, err := h.marketUseCase.GetTicker(stockCode)
				if err == nil && (lastTicker == nil || current.LastPrice != lastTicker.LastPrice || current.Volume != lastTicker.Volume) {
					if err := conn.WriteJSON(map[string]interface{}{"topic": "market.ticker", "data": current}); err != nil {
						return
					}
					lastTicker = current
				}
			}

			if tradeEnabled {
				for _, trade := range h.tradeUseCase.GetTrades(stockCode) {
					if _, seen := seenTradeIDs[trade.ID]; seen {
						continue
					}
					seenTradeIDs[trade.ID] = struct{}{}
					if err := conn.WriteJSON(map[string]interface{}{"topic": "market.trade", "data": trade}); err != nil {
						return
					}
				}
			}
		}
	}
}

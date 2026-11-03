package main

import (
	"log"
	"net/http"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/handler"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
	"github.com/wfpdev/go_mini_exchange/internal/usecase"
)

func main() {
	exchange := &entity.Exchange{
		Stocks:     make(map[string]*entity.Stock),
		OrderBooks: make(map[string]*entity.OrderBook),
		Trades:     make([]*entity.Trade, 0),
	}

	// Sample in-memory stock available to the REST API.
	exchange.AddStock("BBCA", "Bank Central Asia", 10000)
	exchange.AddStock("TLKM", "Telkom Indonesia", 3000)
	exchange.AddStock("TEST", "Test Stock", 100)

	orderRepository := repository.NewInMemoryOrderRepository(exchange)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	orderHandler := handler.NewOrderHandler(createOrderUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("POST /create-order", orderHandler.CreateOrder)

	log.Println("mini exchange listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/wfpdev/go_mini_exchange/internal/entity"
	"github.com/wfpdev/go_mini_exchange/internal/handler"
	"github.com/wfpdev/go_mini_exchange/internal/messaging"
	"github.com/wfpdev/go_mini_exchange/internal/repository"
	"github.com/wfpdev/go_mini_exchange/internal/usecase"
)

func main() {
	exchange := &entity.Exchange{
		Stocks:     make(map[string]*entity.Stock),
		OrderBooks: make(map[string]*entity.OrderBook),
		Trades:     make([]*entity.Trade, 0),
	}

	// Sample in-memory stock available to the system.
	// These are examples only for my learning, not representing anything real.
	exchange.AddStock("BBCA", "Bank Central Asia", 10000)
	exchange.AddStock("TLKM", "Telkom Indonesia", 3000)
	exchange.AddStock("BBRI", "Bank Rakyat Indonesia", 4800)
	exchange.AddStock("BMRI", "Bank Mandiri", 6200)
	exchange.AddStock("ASII", "Astra International", 5100)
	exchange.AddStock("ICBP", "Indofood CBP Sukses Makmur", 11500)
	exchange.AddStock("UNVR", "Unilever Indonesia", 2400)
	exchange.AddStock("GOTO", "GoTo Gojek Tokopedia", 80)
	exchange.AddStock("BBNI", "Bank Negara Indonesia", 4700)
	exchange.AddStock("TEST", "Test Stock", 100)

	orderRepository := repository.NewInMemoryOrderRepository(exchange)

	brokerValue := os.Getenv("KAFKA_BROKERS")
	if brokerValue == "" {
		brokerValue = "localhost:9092"
	}
	brokers := strings.Split(brokerValue, ",")
	topic := os.Getenv("KAFKA_ORDER_TOPIC")
	if topic == "" {
		topic = messaging.OrderTopic
	}

	orderPublisher := messaging.NewKafkaOrderPublisher(brokers, topic)
	defer orderPublisher.Close()

	orderConsumer := messaging.StartKafkaOrderConsumer(brokers, topic, func(orderID int64) error {
		_, err := orderRepository.ProcessPendingOrder(orderID)
		return err
	})
	defer orderConsumer.Close()

	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository, orderPublisher)
	orderHandler := handler.NewOrderHandler(createOrderUseCase)

	tradeRepository := repository.NewInMemoryTradeRepository(exchange)
	tradeUseCase := usecase.NewTradeUseCase(tradeRepository)
	tradeHandler := handler.NewTradeHandler(tradeUseCase)

	marketRepository := repository.NewInMemoryMarketRepository(exchange)
	marketUseCase := usecase.NewMarketUseCase(marketRepository)
	marketHandler := handler.NewMarketHandler(marketUseCase, tradeUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /trades", tradeHandler.GetTrades)
	mux.HandleFunc("GET /ticker", marketHandler.GetTicker)
	mux.HandleFunc("GET /snapshot/ticker", marketHandler.GetTicker)
	mux.HandleFunc("GET /snapshot/order-book", marketHandler.GetOrderBook)
	mux.HandleFunc("GET /ws/market", marketHandler.StreamMarket)
	mux.HandleFunc("POST /create-order", orderHandler.CreateOrder)

	log.Println("mini exchange listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

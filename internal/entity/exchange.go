package entity

import (
	"errors"
	"fmt"
	"time"
)

type Exchange struct {
	Stocks     map[string]*Stock
	OrderBooks map[string]*OrderBook
	Trades     []*Trade

	nextOrderID int64
	nextTradeID int64
}

type OrderBook struct {
	BuyOrders  []*Order
	SellOrders []*Order
}

func (e *Exchange) AddStock(
	code string,
	name string,
	initialPrice float64,
) {
	e.Stocks[code] = &Stock{
		Code:            code,
		Name:            name,
		LastTradedPrice: initialPrice,
	}

	e.OrderBooks[code] = &OrderBook{
		BuyOrders:  []*Order{},
		SellOrders: []*Order{},
	}
}

func (e *Exchange) CreateOrder(
	userID int64,
	stockCode string,
	side OrderSide,
	price float64,
	quantity int,
) (*Order, error) {

	// Validate stock
	if _, exists := e.Stocks[stockCode]; !exists {
		return nil, errors.New("stock not found")
	}

	// Validate order side
	if side != Buy && side != Sell {
		return nil, errors.New("invalid order side")
	}

	// Validate price and quantity
	if price <= 0 || quantity <= 0 {
		return nil, errors.New("price and quantity must be positive")
	}

	now := time.Now()

	if e.nextOrderID == 0 {
		e.nextOrderID++
	}

	order := &Order{
		ID:                e.nextOrderID,
		UserID:            userID,
		StockCode:         stockCode,
		Side:              side,
		Price:             price,
		OriginalQuantity:  quantity,
		FilledQuantity:    0,
		RemainingQuantity: quantity,
		Status:            Pending,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	e.nextOrderID++

	book := e.OrderBooks[stockCode]

	// Match incoming order against existing orders
	e.matchOrder(book, order)

	// Add remaining quantity to the order book
	if order.RemainingQuantity > 0 {
		order.Status = Open
		e.addToOrderBook(book, order)
	}

	return order, nil
}

func (e *Exchange) matchOrder(
	book *OrderBook,
	incoming *Order,
) {

	for incoming.RemainingQuantity > 0 {

		var resting *Order

		// Find the best opposing order
		if incoming.Side == Buy {

			if len(book.SellOrders) == 0 {
				break
			}

			resting = book.SellOrders[0]

			// Buy price must be >= sell price
			if incoming.Price < resting.Price {
				break
			}

		} else {

			if len(book.BuyOrders) == 0 {
				break
			}

			resting = book.BuyOrders[0]

			// Sell price must be <= buy price
			if incoming.Price > resting.Price {
				break
			}
		}

		// Determine executable quantity
		tradeQuantity := incoming.RemainingQuantity

		if resting.RemainingQuantity < tradeQuantity {
			tradeQuantity = resting.RemainingQuantity
		}

		// Trade executes at resting order's price
		tradePrice := resting.Price

		e.executeTrade(
			incoming,
			resting,
			tradePrice,
			tradeQuantity,
		)

		// Update incoming order
		incoming.RemainingQuantity -= tradeQuantity
		incoming.FilledQuantity += tradeQuantity

		// Update resting order
		resting.RemainingQuantity -= tradeQuantity
		resting.FilledQuantity += tradeQuantity

		now := time.Now()

		incoming.UpdatedAt = now
		resting.UpdatedAt = now

		// Update order statuses
		updateOrderStatus(incoming)
		updateOrderStatus(resting)

		// Remove fully filled resting order
		if resting.RemainingQuantity == 0 {

			if incoming.Side == Buy {
				book.SellOrders = book.SellOrders[1:]
			} else {
				book.BuyOrders = book.BuyOrders[1:]
			}
		}
	}
}

func updateOrderStatus(order *Order) {

	if order.RemainingQuantity == 0 {
		order.Status = Filled
		return
	}

	if order.FilledQuantity > 0 {
		order.Status = PartiallyFilled
		return
	}

	order.Status = Open
}

func (e *Exchange) executeTrade(
	incoming *Order,
	resting *Order,
	price float64,
	quantity int,
) {

	var buyOrder *Order
	var sellOrder *Order

	if incoming.Side == Buy {
		buyOrder = incoming
		sellOrder = resting
	} else {
		buyOrder = resting
		sellOrder = incoming
	}

	if e.nextTradeID == 0 {
		e.nextTradeID++
	}

	trade := &Trade{
		ID:          e.nextTradeID,
		StockCode:   incoming.StockCode,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		Price:       price,
		Quantity:    quantity,
		ExecutedAt:  time.Now(),
	}

	e.nextTradeID++

	e.Trades = append(e.Trades, trade)

	// Update last traded price
	stock := e.Stocks[incoming.StockCode]
	stock.LastTradedPrice = price

	fmt.Printf(
		"TRADE: %s | Price: %.2f | Quantity: %d | Buy Order: %d | Sell Order: %d\n",
		trade.StockCode,
		trade.Price,
		trade.Quantity,
		trade.BuyOrderID,
		trade.SellOrderID,
	)
}

func (e *Exchange) addToOrderBook(
	book *OrderBook,
	order *Order,
) {

	if order.Side == Buy {

		book.BuyOrders = append(book.BuyOrders, order)

		// Higher price first
		for i := len(book.BuyOrders) - 1; i > 0; i-- {

			current := book.BuyOrders[i]
			previous := book.BuyOrders[i-1]

			if current.Price > previous.Price {

				book.BuyOrders[i], book.BuyOrders[i-1] =
					book.BuyOrders[i-1], book.BuyOrders[i]

			} else {
				break
			}
		}

	} else {

		book.SellOrders = append(book.SellOrders, order)

		// Lower price first
		for i := len(book.SellOrders) - 1; i > 0; i-- {

			current := book.SellOrders[i]
			previous := book.SellOrders[i-1]

			if current.Price < previous.Price {

				book.SellOrders[i], book.SellOrders[i-1] =
					book.SellOrders[i-1], book.SellOrders[i]

			} else {
				break
			}
		}
	}
}

func (e *Exchange) PrintOrderBook(stockCode string) {

	book := e.OrderBooks[stockCode]
	stock := e.Stocks[stockCode]

	fmt.Println("\n====================")
	fmt.Println("ORDER BOOK:", stockCode)
	fmt.Println("====================")

	fmt.Println("\nSELL ORDERS")

	for _, order := range book.SellOrders {
		fmt.Printf(
			"ID: %d | Price: %.2f | Remaining: %d | Status: %s\n",
			order.ID,
			order.Price,
			order.RemainingQuantity,
			order.Status,
		)
	}

	fmt.Println("\nBUY ORDERS")

	for _, order := range book.BuyOrders {
		fmt.Printf(
			"ID: %d | Price: %.2f | Remaining: %d | Status: %s\n",
			order.ID,
			order.Price,
			order.RemainingQuantity,
			order.Status,
		)
	}

	fmt.Printf(
		"\nLast traded price: %.2f\n",
		stock.LastTradedPrice,
	)
}

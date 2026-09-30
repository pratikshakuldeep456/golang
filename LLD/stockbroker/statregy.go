package stockbroker

// matching_strategy.go

// Contains the matching logic.

// Finds the opposite order.
// Checks whether price can match.
// Creates a Trade.
// Updates quantities.
// Updates order status.
// Keeps remaining LIMIT orders in the book.
type MatchingStrategy interface {
	Match(order *Order, book *OrderBook) []Trade
}

type PriceTimeStrategy struct{}

func (s *PriceTimeStrategy) Match(
	order *Order,
	book *OrderBook,
) []Trade {

	var trades []Trade

	for order.Remaining > 0 {

		opposite := getOppositeOrder(order, book)

		if opposite == nil || !canMatch(order, opposite) {
			break
		}

		trade := executeTrade(order, opposite)
		trades = append(trades, trade)

		updateQuantities(order, opposite)

		if opposite.Remaining == 0 {
			opposite.Status = EXECUTED
			book.Remove(opposite)
		} else {
			opposite.Status = PARTIAL
		}
	}

	updateOrderStatus(order)
	handleRemainingOrder(order, book)

	return trades
}

func getOppositeOrder(order *Order, book *OrderBook) *Order {
	if order.Side == BUY {
		return book.BestSell()
	}
	return book.BestBuy()
}

func canMatch(order *Order, opposite *Order) bool {

	if order.Type == MARKET {
		return true
	}

	if order.Side == BUY {
		return order.Price >= opposite.Price
	}

	return order.Price <= opposite.Price
}

func executeTrade(order *Order, opposite *Order) Trade {

	qty := order.Remaining

	if opposite.Remaining < qty {
		qty = opposite.Remaining
	}

	trade := Trade{
		StockID:  order.StockID,
		Price:    opposite.Price,
		Quantity: qty,
	}

	if order.Side == BUY {
		trade.BuyOrderID = order.ID
		trade.SellOrderID = opposite.ID
	} else {
		trade.BuyOrderID = opposite.ID
		trade.SellOrderID = order.ID
	}

	return trade
}

func updateQuantities(order *Order, opposite *Order) {

	qty := order.Remaining

	if opposite.Remaining < qty {
		qty = opposite.Remaining
	}

	order.Remaining -= qty
	opposite.Remaining -= qty
}

func updateOrderStatus(order *Order) {

	if order.Remaining == 0 {
		order.Status = EXECUTED
	} else if order.Remaining < order.Quantity {
		order.Status = PARTIAL
	} else {
		order.Status = PENDING
	}
}

func handleRemainingOrder(order *Order, book *OrderBook) {

	if order.Remaining == 0 {
		return
	}

	if order.Type == LIMIT {
		book.Add(order)
		return
	}

	order.Status = CANCELLED
}

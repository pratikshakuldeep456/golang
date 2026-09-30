package stockbroker

// . order_book.go

// Stores open orders.

// Keeps BUY orders and SELL orders.
// BestBuy() → highest-priority BUY.
// BestSell() → lowest-priority SELL.
// Add() → adds unmatched order.
// Remove() → removes fully executed order.

// Think: OrderBook = waiting orders.

type OrderBook struct {
	buys  []*Order
	sells []*Order
}

func (b *OrderBook) Add(order *Order) {
	if order.Side == BUY {
		b.buys = append(b.buys, order)
	} else {
		b.sells = append(b.sells, order)
	}
}

func (b *OrderBook) BestBuy() *Order {
	if len(b.buys) == 0 {
		return nil
	}
	return b.buys[0]
}

func (b *OrderBook) BestSell() *Order {
	if len(b.sells) == 0 {
		return nil
	}
	return b.sells[0]
}

func (b *OrderBook) Remove(order *Order) {
	var orders []*Order

	if order.Side == BUY {
		orders = b.buys
	} else {
		orders = b.sells
	}

	for i, o := range orders {
		if o.ID == order.ID {
			orders = append(orders[:i], orders[i+1:]...)
			break
		}
	}

	if order.Side == BUY {
		b.buys = orders
	} else {
		b.sells = orders
	}
}

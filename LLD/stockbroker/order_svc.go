package stockbroker

// . order_service.go
// Entry point for placing an order. orchestatior
// Creates the Order.
// Sends it to MatchingEngine.

type OrderService struct {
	engine *MatchingEngine
}

func NewOrderService(engine *MatchingEngine) *OrderService {
	return &OrderService{
		engine: engine,
	}
}

func (s *OrderService) PlaceOrder(
	id string,
	userID string,
	stockID string,
	side Side,
	orderType OrderType,
	quantity int,
	price float64,
) []Trade {

	order := NewOrder(
		id,
		userID,
		stockID,
		side,
		orderType,
		quantity,
		price,
	)

	return s.engine.PlaceOrder(order)
}

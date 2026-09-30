package stockbroker

// matching_engine.go

// Coordinates matching.

// Maintains an OrderBook for each stock.
// Gets the correct book using StockID.
// Passes the order to the matching strategy.

// Think: Engine = orchestrator.
type MatchingEngine struct {
	books    map[string]*OrderBook
	strategy MatchingStrategy
}

func NewMatchingEngine(strategy MatchingStrategy) *MatchingEngine {
	return &MatchingEngine{
		books:    make(map[string]*OrderBook),
		strategy: strategy,
	}
}

func (e *MatchingEngine) PlaceOrder(order *Order) []Trade {

	book := e.books[order.StockID]

	if book == nil {
		book = &OrderBook{}
		e.books[order.StockID] = book
	}

	return e.strategy.Match(order, book)
}

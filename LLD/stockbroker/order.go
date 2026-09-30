package stockbroker

import "time"

// ==================== ENUMS ====================

// 1. order.go

// Defines the data.

// Order → represents a buy/sell order.
// Trade → represents a successful execution.
// NewOrder() → creates a new order with initial values.
type Side string

const (
	BUY  Side = "BUY"
	SELL Side = "SELL"
)

type OrderType string

const (
	MARKET OrderType = "MARKET"
	LIMIT  OrderType = "LIMIT"
)

type OrderStatus string

const (
	PENDING   OrderStatus = "PENDING"
	PARTIAL   OrderStatus = "PARTIAL"
	EXECUTED  OrderStatus = "EXECUTED"
	CANCELLED OrderStatus = "CANCELLED"
)

// ==================== ORDER ====================

type Order struct {
	ID        string
	UserID    string
	StockID   string
	Side      Side
	Type      OrderType
	Price     float64
	Quantity  int
	Remaining int
	Status    OrderStatus
	Time      time.Time
}

func NewOrder(
	id string,
	userID string,
	stockID string,
	side Side,
	orderType OrderType,
	quantity int,
	price float64,
) *Order {

	return &Order{
		ID:        id,
		UserID:    userID,
		StockID:   stockID,
		Side:      side,
		Type:      orderType,
		Price:     price,
		Quantity:  quantity,
		Remaining: quantity,
		Status:    PENDING,
		Time:      time.Now(),
	}
}

// ==================== TRADE ====================

type Trade struct {
	BuyOrderID  string
	SellOrderID string
	StockID     string
	Price       float64
	Quantity    int
}

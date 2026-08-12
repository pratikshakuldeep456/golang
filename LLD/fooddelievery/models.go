package fooddelivery

type User struct {
	ID      int
	Email   string
	Address string
}

type Cart struct {
	ID     int
	UserId int
	RId    int
	Price  float32
	Items  []CartItem
}

type CartItem struct {
	ID       int
	Name     string
	Quantity int
	Price    float32
}

type Restaurent struct {
	ID       int
	Name     string
	Location string
	About    string
}

type Menu struct {
	ID    int
	RId   int
	Items []MenuItem
}

type MenuItem struct {
	ID       int
	MenuID   int
	FoodMenu string
	Price    float32
}

// --- New: enums needed to complete checkout/payment/order flow ---

// type PaymentStatus string

// const (
// 	PAYMENT_PENDING PaymentStatus = "pending"
// 	PAYMENT_SUCCESS PaymentStatus = "success"
// 	PAYMENT_FAILED  PaymentStatus = "failed"
// )

type Payment struct {
	ID      int
	OrderID int
	Amount  float32
	Status  PaymentStatus
}

// type OrderStatus string

// const (
// 	ORDER_PLACED         OrderStatus = "placed"
// 	ORDER_PAYMENT_FAILED OrderStatus = "payment_failed"
// 	ORDER_CANCELLED      OrderStatus = "cancelled"
// )

type Order struct {
	ID         int
	UserID     int
	ResID      int
	Amount     float32
	Address    string
	Status     OrderStatus
	OrderItems []OrderItem
}

type OrderItem struct {
	ID         int
	MenuItemId int
	Name       string
	Price      int
	Quantity   int
}

// --- New: payment gateway abstraction (3rd party, per scope decision) ---

type PaymentGateway interface {
	Pay(amount float32, orderID int) (string, error)
}

// Stub implementation — in reality this would call Stripe/Razorpay/etc.
type ThirdPartyPaymentGateway struct{}

func (t *ThirdPartyPaymentGateway) Pay(amount float32, orderID int) (string, error) {
	return "payment successful", nil
}

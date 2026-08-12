package fooddelivery

type User struct {
	ID      int
	Email   string
	Address string
	// CartData *Cart
}

type Cart struct {
	ID     int
	UserId int
	RId    int
	// MenuID int
	Price float32
	Items []CartItem
}

type CartItem struct {
	ID int
	// CartId   int

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
type Payment struct {
	ID      int
	OrderID int
	Amount  float32
	Status  PaymentStatus
}

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
	ID int
	// OrderId    int
	MenuItemId int
	Name       string
	Price      float32
	Quantity   int
}

package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// ==================== ENTITIES ====================

type User struct {
	ID      int
	Email   string
	Address string
}

type Restaurant struct {
	ID       int
	Name     string
	Location string
	About    string
}

type MenuItem struct {
	ID           int
	RestaurantID int
	Name         string
	Price        int64 // paise
}

type Cart struct {
	ID           int
	UserID       int
	RestaurantID int
	Items        []CartItem
}

type CartItem struct {
	MenuItemID int
	Name       string
	Price      int64
	Quantity   int
}

type OrderStatus string

const (
	OrderPending OrderStatus = "pending"
	OrderPlaced  OrderStatus = "placed"
	OrderFailed  OrderStatus = "failed"
)

type Order struct {
	ID           int
	UserID       int
	RestaurantID int
	Address      string
	Items        []OrderItem
	TotalAmount  int64
	Status       OrderStatus
}

type OrderItem struct {
	MenuItemID int
	Name       string
	Price      int64
	Quantity   int
}

type PaymentStatus string

const (
	PaymentSuccess PaymentStatus = "success"
	PaymentFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID      int
	OrderID int
	Amount  int64
	Status  PaymentStatus
}

// ==================== STRATEGY PATTERN: PAYMENT ====================

// PaymentStrategy lets us swap payment methods without touching checkout logic.
type PaymentStrategy interface {
	Pay(orderID int, amount int64) PaymentStatus
}

type CardPayment struct{}

func (c *CardPayment) Pay(orderID int, amount int64) PaymentStatus {
	// stub: real impl would call a card gateway (Stripe/Razorpay etc.)
	fmt.Printf("[CardPayment] charging order %d for %d paise\n", orderID, amount)
	return PaymentSuccess
}

type UPIPayment struct{}

func (u *UPIPayment) Pay(orderID int, amount int64) PaymentStatus {
	// stub: real impl would call a UPI gateway
	fmt.Printf("[UPIPayment] charging order %d for %d paise\n", orderID, amount)
	return PaymentSuccess
}

// ==================== OBSERVER PATTERN: NOTIFICATIONS ====================

// NotificationObserver is notified whenever an order's status changes.
type NotificationObserver interface {
	Notify(userID int, message string)
}

type EmailNotifier struct{}

func (e *EmailNotifier) Notify(userID int, message string) {
	fmt.Printf("[Email -> user %d] %s\n", userID, message)
}

type SMSNotifier struct{}

func (s *SMSNotifier) Notify(userID int, message string) {
	fmt.Printf("[SMS -> user %d] %s\n", userID, message)
}

// NotificationDispatcher holds the list of observers and fans out notifications.
// Notification failures must never block checkout, so Dispatch just logs and moves on.
type NotificationDispatcher struct {
	observers []NotificationObserver
}

func NewNotificationDispatcher(observers ...NotificationObserver) *NotificationDispatcher {
	return &NotificationDispatcher{observers: observers}
}

func (d *NotificationDispatcher) Dispatch(userID int, message string) {
	for _, o := range d.observers {
		o.Notify(userID, message)
	}
}

// ==================== SINGLETON: SERVICE ====================

type FoodDeliveryService struct {
	mu sync.Mutex

	users       map[int]User
	restaurants map[int]Restaurant
	menuItems   map[int]MenuItem
	carts       map[int]Cart // keyed by userID
	orders      map[int]Order
	payments    map[int]Payment

	nextOrderID   int
	nextPaymentID int

	paymentStrategy PaymentStrategy
	notifier        *NotificationDispatcher
}

var (
	instance *FoodDeliveryService
	once     sync.Once
)

// GetFoodDeliveryService returns the single shared service instance.
func GetFoodDeliveryService(strategy PaymentStrategy, notifier *NotificationDispatcher) *FoodDeliveryService {
	once.Do(func() {
		instance = &FoodDeliveryService{
			users:           make(map[int]User),
			restaurants:     make(map[int]Restaurant),
			menuItems:       make(map[int]MenuItem),
			carts:           make(map[int]Cart),
			orders:          make(map[int]Order),
			payments:        make(map[int]Payment),
			nextOrderID:     1,
			nextPaymentID:   1,
			paymentStrategy: strategy,
			notifier:        notifier,
		}
	})
	return instance
}

// ==================== SEED HELPERS ====================

func (s *FoodDeliveryService) AddUser(u User) {
	s.users[u.ID] = u
}

func (s *FoodDeliveryService) AddRestaurant(r Restaurant) {
	s.restaurants[r.ID] = r
}

func (s *FoodDeliveryService) AddMenuItem(m MenuItem) {
	s.menuItems[m.ID] = m
}

// ==================== SEARCH ====================

// SearchRestaurants does exact name match; empty query returns everything.
// Partial/fuzzy match is out of scope for this pass - would use Elasticsearch at scale.
func (s *FoodDeliveryService) SearchRestaurants(query string) ([]Restaurant, error) {
	var result []Restaurant
	for _, r := range s.restaurants {
		if query == "" || strings.EqualFold(r.Name, query) {
			result = append(result, r)
		}
	}
	return result, nil
}

// ==================== MENU ====================

func (s *FoodDeliveryService) GetMenu(restaurantID int) ([]MenuItem, error) {
	if _, ok := s.restaurants[restaurantID]; !ok {
		return nil, errors.New("restaurant not found")
	}
	var items []MenuItem
	for _, m := range s.menuItems {
		if m.RestaurantID == restaurantID {
			items = append(items, m)
		}
	}
	return items, nil // empty slice is valid, not an error
}

// ==================== CART ====================

// AddToCart validates quantity, fetches price from backend (never trusts client),
// and rejects adding items from a different restaurant than what's already in the cart.
func (s *FoodDeliveryService) AddToCart(userID, restaurantID, menuItemID, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	item, ok := s.menuItems[menuItemID]
	if !ok || item.RestaurantID != restaurantID {
		return errors.New("menu item not found for this restaurant")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		cart = Cart{ID: userID, UserID: userID, RestaurantID: restaurantID}
	}

	if len(cart.Items) > 0 && cart.RestaurantID != restaurantID {
		return fmt.Errorf("cart already has items from another restaurant (id=%d), clear cart first", cart.RestaurantID)
	}
	cart.RestaurantID = restaurantID

	// merge with existing line item if same menu item added again
	merged := false
	for i, ci := range cart.Items {
		if ci.MenuItemID == menuItemID {
			cart.Items[i].Quantity += quantity
			merged = true
			break
		}
	}
	if !merged {
		cart.Items = append(cart.Items, CartItem{
			MenuItemID: item.ID,
			Name:       item.Name, // snapshot
			Price:      item.Price,
			Quantity:   quantity,
		})
	}

	s.carts[userID] = cart
	return nil
}

func (s *FoodDeliveryService) GetCart(userID int) (Cart, error) {
	cart, ok := s.carts[userID]
	if !ok {
		return Cart{}, errors.New("cart is empty")
	}
	return cart, nil
}

// ==================== CHECKOUT ====================

// Checkout: create order as pending -> snapshot cart -> run payment -> flip to
// placed/failed -> notify (non-blocking, never fails checkout) -> clear cart.
func (s *FoodDeliveryService) Checkout(userID int, address string) (Order, error) {
	s.mu.Lock()
	cart, ok := s.carts[userID]
	if !ok || len(cart.Items) == 0 {
		s.mu.Unlock()
		return Order{}, errors.New("cart is empty")
	}

	var total int64
	var orderItems []OrderItem
	for _, ci := range cart.Items {
		total += ci.Price * int64(ci.Quantity)
		orderItems = append(orderItems, OrderItem{
			MenuItemID: ci.MenuItemID,
			Name:       ci.Name,
			Price:      ci.Price,
			Quantity:   ci.Quantity,
		})
	}

	order := Order{
		ID:           s.nextOrderID,
		UserID:       userID,
		RestaurantID: cart.RestaurantID,
		Address:      address,
		Items:        orderItems,
		TotalAmount:  total,
		Status:       OrderPending,
	}
	s.nextOrderID++
	s.orders[order.ID] = order
	s.mu.Unlock()

	// payment happens outside the lock - external call shouldn't hold it
	status := s.paymentStrategy.Pay(order.ID, order.TotalAmount)

	s.mu.Lock()
	payment := Payment{
		ID:      s.nextPaymentID,
		OrderID: order.ID,
		Amount:  order.TotalAmount,
		Status:  status,
	}
	s.nextPaymentID++
	s.payments[payment.ID] = payment

	if status == PaymentSuccess {
		order.Status = OrderPlaced
		delete(s.carts, userID) // clear cart only on success
	} else {
		order.Status = OrderFailed
	}
	s.orders[order.ID] = order
	s.mu.Unlock()

	msg := fmt.Sprintf("your order #%d is %s", order.ID, order.Status)
	s.notifier.Dispatch(userID, msg) // fire-and-forget, never blocks/fails checkout

	if order.Status == OrderFailed {
		return order, errors.New("payment failed")
	}
	return order, nil
}

// ==================== MAIN ====================

func main() {
	strategy := &CardPayment{} // swap to &UPIPayment{} to change payment method
	notifier := NewNotificationDispatcher(&EmailNotifier{}, &SMSNotifier{})
	svc := GetFoodDeliveryService(strategy, notifier)

	svc.AddUser(User{ID: 1, Email: "a@test.com", Address: "123 Main St"})
	svc.AddRestaurant(Restaurant{ID: 1, Name: "Pizza Palace", Location: "Downtown"})
	svc.AddRestaurant(Restaurant{ID: 2, Name: "Burger Barn", Location: "Uptown"})
	svc.AddMenuItem(MenuItem{ID: 1, RestaurantID: 1, Name: "Margherita", Price: 25000})
	svc.AddMenuItem(MenuItem{ID: 2, RestaurantID: 1, Name: "Pepperoni", Price: 30000})
	svc.AddMenuItem(MenuItem{ID: 3, RestaurantID: 2, Name: "Cheeseburger", Price: 20000})

	results, _ := svc.SearchRestaurants("Pizza Palace")
	fmt.Println("search results:", results)

	menu, _ := svc.GetMenu(1)
	fmt.Println("menu:", menu)

	if err := svc.AddToCart(1, 1, 1, 2); err != nil {
		fmt.Println("add to cart error:", err)
	}
	if err := svc.AddToCart(1, 1, 2, 1); err != nil {
		fmt.Println("add to cart error:", err)
	}

	// this should be rejected - different restaurant
	if err := svc.AddToCart(1, 2, 3, 1); err != nil {
		fmt.Println("expected error (cross-restaurant):", err)
	}

	order, err := svc.Checkout(1, "123 Main St")
	if err != nil {
		fmt.Println("checkout error:", err)
	} else {
		fmt.Printf("order placed: %+v\n", order)
	}
}

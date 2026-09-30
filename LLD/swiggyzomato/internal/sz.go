package internal

import (
	"errors"
	"fmt"
	"sync"
)

// ==================== ENTITIES ====================

type Location struct {
	Lat, Lng float64
}

type User struct {
	ID      string
	Name    string
	Address string
}

type MenuItem struct {
	ID           string
	Name         string
	Price        float64
	RestaurantID string
}

type Restaurant struct {
	ID    string
	Name  string
	Items map[string]*MenuItem
}

type CartItem struct {
	ItemID   string
	Quantity int
	Price    float64
}

type Cart struct {
	UserID       string
	RestaurantID string
	Items        map[string]*CartItem
}

type OrderItem struct {
	ItemID   string
	Quantity int
	Price    float64
}

type Order struct {
	ID           string
	UserID       string
	RestaurantID string
	Items        []OrderItem
	Amount       float64
	State        OrderState
	ExecutiveID  string
}

type DeliveryExecutive struct {
	ID        string
	Location  Location
	Available bool
}

// ==================== STATE PATTERN ====================

type OrderState interface {
	Name() string
	Prepare(*Order) error
	Ready(*Order) error
	Dispatch(*Order) error
	Deliver(*Order) error
}

type BaseState struct{}

func (BaseState) Prepare(*Order) error  { return errors.New("invalid transition") }
func (BaseState) Ready(*Order) error    { return errors.New("invalid transition") }
func (BaseState) Dispatch(*Order) error { return errors.New("invalid transition") }
func (BaseState) Deliver(*Order) error  { return errors.New("invalid transition") }

type PlacedState struct{ BaseState }

func (PlacedState) Name() string { return "ORDER_PLACED" }

func (PlacedState) Prepare(o *Order) error {
	o.State = PreparingState{}
	return nil
}

type PreparingState struct{ BaseState }

func (PreparingState) Name() string { return "PREPARING" }

func (PreparingState) Ready(o *Order) error {
	o.State = ReadyState{}
	return nil
}

type ReadyState struct{ BaseState }

func (ReadyState) Name() string { return "READY_FOR_PICKUP" }

func (ReadyState) Dispatch(o *Order) error {
	o.State = OutForDeliveryState{}
	return nil
}

type OutForDeliveryState struct{ BaseState }

func (OutForDeliveryState) Name() string { return "OUT_FOR_DELIVERY" }

func (OutForDeliveryState) Deliver(o *Order) error {
	o.State = DeliveredState{}
	return nil
}

type DeliveredState struct{ BaseState }

func (DeliveredState) Name() string { return "DELIVERED" }

// ==================== PAYMENT ====================

type PaymentProvider interface {
	Pay(orderID string, amount float64) error
}

// Third-party provider adapter
type RazorpayAdapter struct{}

func (RazorpayAdapter) Pay(orderID string, amount float64) error {
	fmt.Printf("Payment successful: order=%s amount=%.2f\n", orderID, amount)
	return nil
}

// ==================== DELIVERY STRATEGY ====================

type AssignmentStrategy interface {
	Assign(*Order, []*DeliveryExecutive) (*DeliveryExecutive, error)
}

type NearestExecutiveStrategy struct{}

func (NearestExecutiveStrategy) Assign(
	order *Order,
	executives []*DeliveryExecutive,
) (*DeliveryExecutive, error) {

	// Simplified for LLD.
	// Real implementation would calculate geo distance.

	for _, e := range executives {
		if e.Available {
			return e, nil
		}
	}

	return nil, errors.New("no delivery executive available")
}

// ==================== EVENT PUBLISHER ====================

type Event struct {
	OrderID string
	State   string
}

type EventPublisher interface {
	Publish(Event)
}

type SimpleEventPublisher struct{}

func (SimpleEventPublisher) Publish(e Event) {
	fmt.Printf("EVENT: order=%s state=%s\n", e.OrderID, e.State)
}

// ==================== SERVICES ====================

type RestaurantService struct {
	restaurants map[string]*Restaurant
}

func (s *RestaurantService) Search(name string) []*Restaurant {
	var result []*Restaurant

	for _, r := range s.restaurants {
		if r.Name == name {
			result = append(result, r)
		}
	}

	return result
}

func (s *RestaurantService) GetMenu(restaurantID string) ([]*MenuItem, error) {
	r, ok := s.restaurants[restaurantID]
	if !ok {
		return nil, errors.New("restaurant not found")
	}

	result := make([]*MenuItem, 0, len(r.Items))

	for _, item := range r.Items {
		result = append(result, item)
	}

	return result, nil
}

// ==================== CART SERVICE ====================

type CartService struct {
	carts map[string]*Cart
	mu    sync.RWMutex
}

func (s *CartService) AddItem(
	userID string,
	item *MenuItem,
	qty int,
) error {

	if qty <= 0 {
		return errors.New("invalid quantity")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]

	if !ok {
		cart = &Cart{
			UserID:       userID,
			RestaurantID: item.RestaurantID,
			Items:        make(map[string]*CartItem),
		}

		s.carts[userID] = cart
	}

	// One restaurant per cart.
	if cart.RestaurantID != item.RestaurantID {
		return errors.New("cart can contain items from one restaurant only")
	}

	if existing, ok := cart.Items[item.ID]; ok {
		existing.Quantity += qty
	} else {
		cart.Items[item.ID] = &CartItem{
			ItemID:   item.ID,
			Quantity: qty,
			Price:    item.Price,
		}
	}

	return nil
}

func (s *CartService) RemoveItem(userID, itemID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]
	if !ok {
		return errors.New("cart not found")
	}

	delete(cart.Items, itemID)
	return nil
}

func (s *CartService) GetCart(userID string) (*Cart, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.carts[userID]
	if !ok {
		return nil, errors.New("cart not found")
	}

	return cart, nil
}

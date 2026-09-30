package internal

import (
	"errors"
	"fmt"
	"sync"
)

// ==================== ORDER SERVICE ====================

type OrderService struct {
	orders     map[string]*Order
	cart       *CartService
	payment    PaymentProvider
	delivery   AssignmentStrategy
	executives []*DeliveryExecutive
	publisher  EventPublisher

	mu     sync.Mutex
	nextID int
}

func (s *OrderService) Checkout(userID string) (*Order, error) {

	// Serialize order creation in this in-memory implementation.
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, err := s.cart.GetCart(userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	order := &Order{
		ID:           fmt.Sprintf("ORD-%d", s.nextID),
		UserID:       userID,
		RestaurantID: cart.RestaurantID,
		State:        PlacedState{},
	}

	s.nextID++

	for _, item := range cart.Items {
		order.Items = append(order.Items, OrderItem{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
			Price:    item.Price,
		})

		order.Amount += item.Price * float64(item.Quantity)
	}

	// Payment is external.
	if err := s.payment.Pay(order.ID, order.Amount); err != nil {
		return nil, err
	}

	// Save order after successful payment.
	s.orders[order.ID] = order

	// Publish state change.
	s.publisher.Publish(Event{
		OrderID: order.ID,
		State:   order.State.Name(),
	})

	// Assign delivery executive.
	executive, err := s.delivery.Assign(order, s.executives)
	if err != nil {
		return nil, err
	}

	executive.Available = false
	order.ExecutiveID = executive.ID

	return order, nil
}

// ==================== ORDER STATE UPDATES ====================

func (s *OrderService) updateState(
	orderID string,
	transition func(OrderState, *Order) error,
) error {

	order, ok := s.orders[orderID]
	if !ok {
		return errors.New("order not found")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	oldState := order.State.Name()

	if err := transition(order.State, order); err != nil {
		return err
	}

	if oldState != order.State.Name() {
		s.publisher.Publish(Event{
			OrderID: order.ID,
			State:   order.State.Name(),
		})
	}

	return nil
}

func (s *OrderService) Prepare(orderID string) error {
	return s.updateState(orderID,
		func(state OrderState, o *Order) error {
			return state.Prepare(o)
		})
}

func (s *OrderService) Ready(orderID string) error {
	return s.updateState(orderID,
		func(state OrderState, o *Order) error {
			return state.Ready(o)
		})
}

func (s *OrderService) Dispatch(orderID string) error {
	return s.updateState(orderID,
		func(state OrderState, o *Order) error {
			return state.Dispatch(o)
		})
}

func (s *OrderService) Deliver(orderID string) error {
	return s.updateState(orderID,
		func(state OrderState, o *Order) error {
			return state.Deliver(o)
		})
}

// ==================== MAIN ====================

func SzTest() {

	item := &MenuItem{
		ID:           "P1",
		Name:         "Pizza",
		Price:        300,
		RestaurantID: "R1",
	}

	restaurant := &Restaurant{
		ID:   "R1",
		Name: "Dominos",
		Items: map[string]*MenuItem{
			item.ID: item,
		},
	}

	restaurantService := &RestaurantService{
		restaurants: map[string]*Restaurant{
			restaurant.ID: restaurant,
		},
	}

	cartService := &CartService{
		carts: make(map[string]*Cart),
	}

	_ = restaurantService

	cartService.AddItem("U1", item, 2)

	deliveryStrategy := NearestExecutiveStrategy{}

	orderService := &OrderService{
		orders:   make(map[string]*Order),
		cart:     cartService,
		payment:  RazorpayAdapter{},
		delivery: deliveryStrategy,
		executives: []*DeliveryExecutive{
			{
				ID:        "D1",
				Available: true,
			},
		},
		publisher: SimpleEventPublisher{},
		nextID:    1,
	}

	order, err := orderService.Checkout("U1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Order created:", order.ID)

	orderService.Prepare(order.ID)
	orderService.Ready(order.ID)
	orderService.Dispatch(order.ID)
	orderService.Deliver(order.ID)

	fmt.Println("Final state:", order.State.Name())
}

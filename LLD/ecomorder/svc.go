package ecomcart

import (
	"errors"
	"fmt"
	"log"

	"sync"
	"time"
)

type OrderStatus string

const (
	Confirmed OrderStatus = "CONFIRMED"
)

type Product struct {
	ID          string
	Name        string
	Price       float64
	MaxPerOrder int // 0 = no limit
}

type Inventory struct {
	ProductID string
	Quantity  int
}

type CartItem struct {
	ProductID string
	Quantity  int
}

type Cart struct {
	ID     string
	UserID string
	Items  map[string]*CartItem
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type Order struct {
	ID     string
	UserID string
	Items  []OrderItem
	Total  float64
	Status OrderStatus
}

type CartService struct {
	mu        sync.Mutex
	products  map[string]*Product
	inventory map[string]*Inventory
	carts     map[string]*Cart // userID -> cart
	orders    map[string]*Order
}

func NewCartService() *CartService {
	return &CartService{
		products:  make(map[string]*Product),
		inventory: make(map[string]*Inventory),
		carts:     make(map[string]*Cart),
		orders:    make(map[string]*Order),
	}
}

// ---------------- Product / Inventory ----------------

func (s *CartService) AddProduct(product *Product, quantity int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.products[product.ID] = product
	s.inventory[product.ID] = &Inventory{
		ProductID: product.ID,
		Quantity:  quantity,
	}
}

// ---------------- Cart APIs ----------------

func (s *CartService) AddToCart(
	userID string,
	productID string,
	qty int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if qty <= 0 {
		return errors.New("invalid quantity")
	}

	product, ok := s.products[productID]
	if !ok {
		return errors.New("product not found")
	}

	cart := s.carts[userID]

	if cart == nil {
		cart = &Cart{
			ID:     userID + "-cart",
			UserID: userID,
			Items:  make(map[string]*CartItem),
		}
		s.carts[userID] = cart
	}

	currentQty := 0

	if item, ok := cart.Items[productID]; ok {
		currentQty = item.Quantity
	}

	newQty := currentQty + qty

	// Per-order purchase limit.
	if product.MaxPerOrder > 0 &&
		newQty > product.MaxPerOrder {
		return errors.New("purchase limit exceeded")
	}

	cart.Items[productID] = &CartItem{
		ProductID: productID,
		Quantity:  newQty,
	}

	return nil
}

func (s *CartService) UpdateQuantity(
	userID string,
	productID string,
	qty int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if qty <= 0 {
		return errors.New("invalid quantity")
	}

	product, ok := s.products[productID]
	if !ok {
		return errors.New("product not found")
	}

	if product.MaxPerOrder > 0 &&
		qty > product.MaxPerOrder {
		return errors.New("purchase limit exceeded")
	}

	cart, ok := s.carts[userID]
	if !ok {
		return errors.New("cart not found")
	}

	item, ok := cart.Items[productID]
	if !ok {
		return errors.New("product not in cart")
	}

	item.Quantity = qty

	return nil
}

func (s *CartService) RemoveFromCart(
	userID string,
	productID string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]
	if !ok {
		return errors.New("cart not found")
	}

	if _, ok := cart.Items[productID]; !ok {
		return errors.New("product not in cart")
	}

	delete(cart.Items, productID)

	return nil
}

func (s *CartService) GetCart(userID string) (*Cart, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]
	if !ok {
		return nil, errors.New("cart not found")
	}

	return cart, nil
}

// ---------------- Checkout ----------------

func (s *CartService) Checkout(userID string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]

	if !ok || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// First validate EVERYTHING.
	// We don't modify anything until validation succeeds.
	order := &Order{
		ID:     generateID(),
		UserID: userID,
		Status: Confirmed,
	}

	for _, cartItem := range cart.Items {

		product, ok := s.products[cartItem.ProductID]
		if !ok {
			return nil, errors.New("product not found")
		}

		inventory, ok := s.inventory[cartItem.ProductID]
		if !ok {
			return nil, errors.New("inventory not found")
		}

		// Re-check purchase limit at checkout.
		if product.MaxPerOrder > 0 &&
			cartItem.Quantity > product.MaxPerOrder {
			return nil, errors.New("purchase limit exceeded")
		}

		// Check inventory.
		if inventory.Quantity < cartItem.Quantity {
			return nil, fmt.Errorf(
				"insufficient inventory for product %s",
				product.ID,
			)
		}

		order.Items = append(order.Items, OrderItem{
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
		})

		order.Total += product.Price * float64(cartItem.Quantity)
	}

	// All validation succeeded.
	// Now perform mutations.

	for _, item := range order.Items {
		s.inventory[item.ProductID].Quantity -= item.Quantity
	}

	s.orders[order.ID] = order

	// Clear cart only after successful checkout.
	cart.Items = make(map[string]*CartItem)

	return order, nil
}

func generateID() string {
	return fmt.Sprintf("order-%d", time.Now().UnixNano())
}

func testCodeEcom() {
	service := NewCartService()

	// Add products + inventory
	service.AddProduct(&Product{
		ID:          "p1",
		Name:        "iPhone",
		Price:       1000,
		MaxPerOrder: 2,
	}, 10)

	service.AddProduct(&Product{
		ID:          "p2",
		Name:        "T-Shirt",
		Price:       50,
		MaxPerOrder: 5,
	}, 20)

	// Add to cart
	if err := service.AddToCart("user1", "p1", 2); err != nil {
		log.Println(err)
	}

	if err := service.AddToCart("user1", "p2", 3); err != nil {
		log.Println(err)
	}

	// View cart
	cart, err := service.GetCart("user1")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Cart:", cart)

	// Update quantity
	if err := service.UpdateQuantity("user1", "p2", 2); err != nil {
		log.Println(err)
	}

	// Remove product
	if err := service.RemoveFromCart("user1", "p2"); err != nil {
		log.Println(err)
	}

	// Checkout
	order, err := service.Checkout("user1")
	if err != nil {
		log.Println("Checkout failed:", err)
		return
	}

	fmt.Println("Order ID:", order.ID)
	fmt.Println("Total:", order.Total)
	fmt.Println("Status:", order.Status)
}

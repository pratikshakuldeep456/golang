// package service

// import (
// 	"fmt"
// 	fooddelivery "pratikshakuldeep456/golang/LLD/foodDelievery"
// 	"sync"
// )

// type FoodDelieverySVC struct {
// 	Restaurants     map[int]*fooddelivery.Restaurent
// 	RestaurantsMenu map[int][]fooddelivery.Menu
// 	Users           map[int]*fooddelivery.User
// 	Carts           map[int]*fooddelivery.Cart
// 	Orders          map[int][]fooddelivery.Order
// 	Mu              sync.Mutex
// }

// var once sync.Once
// var FDInstanceSvc *FoodDelieverySVC

// func GetInstance() *FoodDelieverySVC {
// 	once.Do(
// 		func() {
// 			FDInstanceSvc = &FoodDelieverySVC{
// 				Restaurants:     make(map[int]*fooddelivery.Restaurent),
// 				RestaurantsMenu: make(map[int][]fooddelivery.Menu),
// 				Users:           make(map[int]*fooddelivery.User),
// 				Carts:           make(map[int]*fooddelivery.Cart),
// 				Orders:          make(map[int][]fooddelivery.Order),
// 			}

// 		})
// 	return FDInstanceSvc
// }

// func (fd *FoodDelieverySVC) ListResraurent() /*[]fooddelivery.Restaurent*/ {
// 	if len(fd.Restaurants) == 0 {
// 		fmt.Println("no restaurent found")
// 	}
// 	// var result *[]fooddelivery.Restaurent
// 	for i, r := range fd.Restaurants {
// 		//fmt.Println("resId:", i, "name:", r.Name, "location:", r.Location)
// 		//aka
// 		fmt.Printf("ResId: %d, Name: %s, Location: %s\n", i, r.Name, r.Location)
// 	}

// }

// func (fd *FoodDelieverySVC) AddRestarent(r *fooddelivery.Restaurent) (string, error) {

// 	if r == nil {
// 		return "", fmt.Errorf(" restaurnet cant be empty")
// 	}
// 	fd.Mu.Lock()
// 	defer fd.Mu.Unlock()
// 	fd.Restaurants[r.ID] = r
// 	return "res added", nil
// }

// func (fd *FoodDelieverySVC) ViewMenu(id int) []fooddelivery.Menu {
// 	if len(fd.RestaurantsMenu[id]) == 0 {
// 		return []fooddelivery.Menu{}
// 	}
// 	menuList := fd.RestaurantsMenu[id]
// 	if len(menuList) == 0 {
// 		return []fooddelivery.Menu{}
// 	}

// 	for _, j := range menuList {
// 		fmt.Println(" ", j.ID, j.FoodMenu, j.Price, j.RId)
// 	}
// 	return fd.RestaurantsMenu[id]
// }

// func (fd *FoodDelieverySVC) AddMenu(menu *fooddelivery.Menu) {
// 	if menu == nil {
// 		fmt.Errorf("food menu cant be emptpy")
// 	}

// 	fd.RestaurantsMenu[menu.RId] = append(fd.RestaurantsMenu[menu.RId], *menu)
// 	fmt.Println(" hey menu added")

// }

// func (fd *FoodDelieverySVC) AddtoCart(mid int, userid, resid int, quantity int, price float32) (string, error) {
// 	fd.Mu.Lock()
// 	defer fd.Mu.Unlock()
// 	cart1 := &fooddelivery.Cart{
// 		ID:     1,
// 		UserId: userid,
// 		RId:    resid,
// 		Price:  price,
// 	}
// 	fd.Carts[cart1.ID] = cart1
// 	return "item added to the caart", nil

// }

package service

import (
	"fmt"
	"sync"

	fooddelivery "pratikshakuldeep456/golang/LLD/foodDelievery"
	"pratikshakuldeep456/golang/LLD/fooddelievery/service"
)

type FoodDelieverySVC struct {
	Restaurants     map[int]*fooddelivery.Restaurent
	RestaurantsMenu map[int][]fooddelivery.Menu
	Users           map[int]*fooddelivery.User
	Carts           map[int]*fooddelivery.Cart
	Orders          map[int][]fooddelivery.Order
	Payments        map[int]*fooddelivery.Payment
	Mu              sync.Mutex
}

var once sync.Once
var FDInstanceSvc *FoodDelieverySVC

func GetInstance() *FoodDelieverySVC {
	once.Do(func() {
		FDInstanceSvc = &FoodDelieverySVC{
			Restaurants:     make(map[int]*fooddelivery.Restaurent),
			RestaurantsMenu: make(map[int][]fooddelivery.Menu),
			Users:           make(map[int]*fooddelivery.User),
			Carts:           make(map[int]*fooddelivery.Cart),
			Orders:          make(map[int][]fooddelivery.Order),
			Payments:        make(map[int]*fooddelivery.Payment),
		}
	})
	return FDInstanceSvc
}

// ---------- User ----------

func (fd *FoodDelieverySVC) RegisterUser(email, address string) (*fooddelivery.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	user := &fooddelivery.User{
		ID:      generateUserID(),
		Email:   email,
		Address: address,
	}
	fd.Users[user.ID] = user
	return user, nil
}

// ---------- Restaurant ----------

func (fd *FoodDelieverySVC) AddRestarent(r *fooddelivery.Restaurent) (string, error) {
	if r == nil {
		return "", fmt.Errorf("restaurant cannot be empty")
	}
	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	if r.ID == 0 {
		r.ID = generateRestaurantID()
	}
	fd.Restaurants[r.ID] = r
	return "restaurant added", nil
}

func (fd *FoodDelieverySVC) ListResraurent() []*fooddelivery.Restaurent {
	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	if len(fd.Restaurants) == 0 {
		fmt.Println("no restaurants found")
		return []*fooddelivery.Restaurent{}
	}

	result := make([]*fooddelivery.Restaurent, 0, len(fd.Restaurants))
	for id, r := range fd.Restaurants {
		fmt.Printf("ResId: %d, Name: %s, Location: %s\n", id, r.Name, r.Location)
		result = append(result, r)
	}
	return result
}

// ---------- Menu ----------

func (fd *FoodDelieverySVC) AddMenu(menu *fooddelivery.Menu) (string, error) {
	if menu == nil {
		return "", fmt.Errorf("menu cannot be empty")
	}
	if _, exists := fd.Restaurants[menu.RId]; !exists {
		return "", fmt.Errorf("restaurant %d not found", menu.RId)
	}

	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	fd.RestaurantsMenu[menu.RId] = append(fd.RestaurantsMenu[menu.RId], *menu)
	return "menu added", nil
}

func (fd *FoodDelieverySVC) ViewMenu(id int) []fooddelivery.Menu {
	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	menuList, exists := fd.RestaurantsMenu[id]
	if !exists || len(menuList) == 0 {
		fmt.Println("no menu found for restaurant", id)
		return []fooddelivery.Menu{}
	}

	for _, m := range menuList {
		for _, item := range m.Items {
			fmt.Println(" ", item.ID, item.FoodMenu, item.Price)
		}
	}
	return menuList
}

// ---------- Cart ----------

func (fd *FoodDelieverySVC) AddtoCart(userid, resid int, name string, quantity int, price float32) (string, error) {
	if quantity <= 0 {
		return "", fmt.Errorf("quantity must be positive")
	}
	if _, exists := fd.Users[userid]; !exists {
		return "", fmt.Errorf("user %d not found", userid)
	}
	if _, exists := fd.Restaurants[resid]; !exists {
		return "", fmt.Errorf("restaurant %d not found", resid)
	}

	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	cart, exists := fd.Carts[userid]
	if !exists {
		cart = &fooddelivery.Cart{
			ID:     userid,
			UserId: userid,
			RId:    resid,
			Items:  []fooddelivery.CartItem{},
		}
		fd.Carts[userid] = cart
	}

	// guard: don't silently mix items from two different restaurants in one cart
	if cart.RId != resid && len(cart.Items) > 0 {
		return "", fmt.Errorf("cart already has items from a different restaurant; checkout or clear cart first")
	}

	cart.Items = append(cart.Items, fooddelivery.CartItem{
		ID:       generateCartItemID(),
		Name:     name,
		Quantity: quantity,
		Price:    price,
	})
	cart.Price += price * float32(quantity)

	return "item added to cart", nil
}

// ---------- Checkout ----------

func (fd *FoodDelieverySVC) Checkout(userid int, gateway fooddelivery.PaymentGateway) (*fooddelivery.Order, error) {
	fd.Mu.Lock()
	defer fd.Mu.Unlock()

	cart, exists := fd.Carts[userid]
	if !exists || len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	user, exists := fd.Users[userid]
	if !exists {
		return nil, fmt.Errorf("user %d not found", userid)
	}

	var total float32
	orderItems := make([]fooddelivery.OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		total += item.Price * float32(item.Quantity)
		orderItems = append(orderItems, fooddelivery.OrderItem{
			ID:       generateOrderItemID(),
			Name:     item.Name,
			Price:    int(item.Price),
			Quantity: item.Quantity,
		})
	}

	order := &fooddelivery.Order{
		ID:         generateOrderID(),
		UserID:     userid,
		ResID:      cart.RId,
		Amount:     total,
		Address:    user.Address,
		Status:     service.ORDER_PAYMENT_FAILED, // pessimistic default until payment confirms
		OrderItems: orderItems,
	}

	_, err := gateway.Pay(total, order.ID)

	payment := &fooddelivery.Payment{
		ID:      generatePaymentID(),
		OrderID: order.ID,
		Amount:  total,
	}

	if err != nil {
		payment.Status = service.PAYMENT_FAILED
		fd.Payments[payment.ID] = payment
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	payment.Status = service.PAYMENT_SUCCESS
	fd.Payments[payment.ID] = payment
	order.Status = service.ORDER_PLACED

	fd.Orders[userid] = append(fd.Orders[userid], *order)
	delete(fd.Carts, userid) // clear cart only after successful payment

	fd.NotifyOrderPlaced(order)

	return order, nil
}

// ---------- Notification (Observer-style hook) ----------

func (fd *FoodDelieverySVC) NotifyOrderPlaced(order *fooddelivery.Order) {
	fmt.Printf("Notification: order #%d placed for user %d, total: %.2f\n",
		order.ID, order.UserID, order.Amount)
}

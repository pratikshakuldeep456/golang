package service

import (
	"fmt"
	fooddelivery "pratikshakuldeep456/golang/LLD/foodDelievery"
	"sync"
)

type FoodDelieverySVC struct {
	Restaurants     map[int]*fooddelivery.Restaurent
	RestaurantsMenu map[int][]fooddelivery.Menu
	Users           map[int]*fooddelivery.User
	Carts           map[int]*fooddelivery.Cart
	Orders          map[int][]fooddelivery.Order
	Mu              sync.Mutex
}

var once sync.Once
var FDInstanceSvc *FoodDelieverySVC

func GetInstance() *FoodDelieverySVC {
	once.Do(
		func() {
			FDInstanceSvc = &FoodDelieverySVC{
				Restaurants:     make(map[int]*fooddelivery.Restaurent),
				RestaurantsMenu: make(map[int][]fooddelivery.Menu),
				Users:           make(map[int]*fooddelivery.User),
				Carts:           make(map[int]*fooddelivery.Cart),
				Orders:          make(map[int][]fooddelivery.Order),
			}

		})
	return FDInstanceSvc
}

func (fd *FoodDelieverySVC) ListResraurent() /*[]fooddelivery.Restaurent*/ {
	if len(fd.Restaurants) == 0 {
		fmt.Println("no restaurent found")
	}
	// var result *[]fooddelivery.Restaurent
	for i, r := range fd.Restaurants {
		//fmt.Println("resId:", i, "name:", r.Name, "location:", r.Location)
		//aka
		fmt.Printf("ResId: %d, Name: %s, Location: %s\n", i, r.Name, r.Location)
	}

}

func (fd *FoodDelieverySVC) AddRestarent(r *fooddelivery.Restaurent) (string, error) {

	if r == nil {
		return "", fmt.Errorf(" restaurnet cant be empty")
	}
	fd.Mu.Lock()
	defer fd.Mu.Unlock()
	fd.Restaurants[r.ID] = r
	return "res added", nil
}

func (fd *FoodDelieverySVC) ViewMenu(id int) []fooddelivery.Menu {
	if len(fd.RestaurantsMenu[id]) == 0 {
		return []fooddelivery.Menu{}
	}
	menuList := fd.RestaurantsMenu[id]
	if len(menuList) == 0 {
		return []fooddelivery.Menu{}
	}

	for _, j := range menuList {
		fmt.Println(" ", j.ID, j.FoodMenu, j.Price, j.RId)
	}
	return fd.RestaurantsMenu[id]
}

func (fd *FoodDelieverySVC) AddMenu(menu *fooddelivery.Menu) {
	if menu == nil {
		fmt.Errorf("food menu cant be emptpy")
	}

	fd.RestaurantsMenu[menu.RId] = append(fd.RestaurantsMenu[menu.RId], *menu)
	fmt.Println(" hey menu added")

}

func (fd *FoodDelieverySVC) AddtoCart(mid int, userid, resid int, quantity int, price float32) (string, error) {
	fd.Mu.Lock()
	defer fd.Mu.Unlock()
	cart1 := &fooddelivery.Cart{
		ID:     1,
		UserId: userid,
		RId:    resid,
		Price:  price,
	}
	fd.Carts[cart1.ID] = cart1
	return "item added to the caart", nil

}

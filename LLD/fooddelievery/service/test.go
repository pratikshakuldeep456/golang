package service

import (
	"fmt"

	fooddelivery "pratikshakuldeep456/golang/LLD/foodDelievery"
)

func FDTest() {
	fd := GetInstance()
	gateway := &fooddelivery.ThirdPartyPaymentGateway{}

	fmt.Println("=== Register user ===")
	user, err := fd.RegisterUser("piku@example.com", "123 MG Road")
	check("RegisterUser", err)

	fmt.Println("\n=== Add restaurant ===")
	_, err = fd.AddRestarent(&fooddelivery.Restaurent{Name: "Punjabi Tadka", Location: "Koramangala"})
	check("AddRestarent", err)

	var resID int
	for id := range fd.Restaurants {
		resID = id
	}

	fmt.Println("\n=== Add menu ===")
	_, err = fd.AddMenu(&fooddelivery.Menu{
		RId: resID,
		Items: []fooddelivery.MenuItem{
			{ID: 1, FoodMenu: "Butter Naan", Price: 40},
			{ID: 2, FoodMenu: "Paneer Tikka", Price: 180},
		},
	})
	check("AddMenu", err)

	fmt.Println("\n=== View menu ===")
	fd.ViewMenu(resID)

	fmt.Println("\n=== Add to cart ===")
	_, err = fd.AddtoCart(user.ID, resID, "Butter Naan", 2, 40)
	check("AddtoCart item 1", err)
	_, err = fd.AddtoCart(user.ID, resID, "Paneer Tikka", 1, 180)
	check("AddtoCart item 2", err)

	fmt.Println("\n=== Checkout ===")
	order, err := fd.Checkout(user.ID, gateway)
	check("Checkout", err)
	if order != nil {
		fmt.Printf("Order placed: ID=%d, Amount=%.2f, Status=%s\n", order.ID, order.Amount, order.Status)
	}

	fmt.Println("\n=== Checkout again with empty cart — should fail ===")
	_, err = fd.Checkout(user.ID, gateway)
	check("Checkout (empty cart, expect error)", err)

	fmt.Println("\n=== Add to cart — invalid user ===")
	_, err = fd.AddtoCart(9999, resID, "Butter Naan", 1, 40)
	check("AddtoCart (invalid user, expect error)", err)
}

func check(label string, err error) {
	if err != nil {
		fmt.Println("FAILED:", label, "-", err)
	} else {
		fmt.Println("OK:", label)
	}
}

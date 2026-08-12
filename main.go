package main

import (
	"fmt"
	"pratikshakuldeep456/golang/LLD/ridesharingsystem"
)

func Counter() func() int {
	count := 0
	fmt.Println("Counter initialized", count)
	return func() int {
		count++

		return count
	}

}

func main() {
	// notes.Main2()
	// notes.Context()
	// notes.NonBlockingCh()
	// notes.WorkerPool()
	// notes.WaitGrp()
	// notes.Limit()
	// notes.Mutex()
	// notes.FormatString()
	// notes.Json()
	// notes.Parse()
	// notes.ReadWrite()
	// notes.ContextFunc()

	// cmd := exec.Command("echo", "Hello from child process")
	// cmd.Run() // Spawns the process
	// User1 := solid.User{Name: "Pratiksha", Age: 20}
	// data, err := User1.CheckCategory()
	// fmt.Println(data, err)

	// // open close principle
	// rectangle := solid.Rectange{LENGTH: 10, BREADTH: 5}
	// circle := solid.Circle{RADIUS: 7}

	// fmt.Println("Area of Rectangle:", solid.GetArea(&rectangle))
	// fmt.Println("Area of Circle:", solid.GetArea(&circle))

	// //lsv

	// svAccount := solid.SavingsAccount{Balance: 1000}
	// svAccount.Deposit(500)
	// fmt.Println("Savings Account Balance:", svAccount.GetBalance())

	// svAccount.Withdraw(200)
	// fmt.Println("Savings Account Balance after withdrawal:", svAccount.GetBalance())

	// fdAccount := solid.FDAccount{Balance: 2000}
	// fdAccount.Deposit(1000)
	// fmt.Println("FD Account Balance:", fdAccount.GetBalance())

	// // design patterns
	// emaila := designpattern.EmaiL{To: "sd", Subject: "Hello", Body: "This is a test email."}
	// notifcationService := designpattern.NewNotificationService(&emaila)
	// err = notifcationService.Notify()
	// if err != nil {
	// 	fmt.Println("Error sending notification:", err)
	// } else {
	// 	fmt.Println("Notification sent successfully")

	// }

	// upiPayment := designpattern.UPI{UpiID: "pratiksha@upi", Amount: 100.0}

	// paymentSystem := designpattern.NewPaymentSystem("upi", &designpattern.PaymentConfig{Amount: upiPayment.Amount})
	// err = paymentSystem.ProcessPayment(upiPayment.Amount)
	// if err != nil {
	// 	fmt.Println("Error processing payment:", err)
	// } else {
	// 	fmt.Println("Payment processed successfully")
	// }

	// stock := &designpattern.StockPrice{Symbol: "AAPL", Price: 150}

	// stock.Attach(&designpattern.SMSNotifier{PhoneNumber: "+1234567890"})
	// stock.Attach(&designpattern.EmailNotifier{EmailAddress: "trader@example.com"})

	// stock.Detach(&designpattern.SMSNotifier{PhoneNumber: "+1234567890"})

	// fooddelivery.User

	// svc := service.GetInstance()
	// svc.ListResraurent()
	// res1 := &fooddelivery.Restaurent{
	// 	ID:       1,
	// 	Name:     "meghana",
	// 	Location: "hsr",
	// 	About:    "nonveg sjdfdgjhjkgfsdk",
	// }
	// res2 := &fooddelivery.Restaurent{
	// 	ID:       2,
	// 	Name:     "meghana",
	// 	Location: "bellendur",
	// 	About:    "nonveg sjdfdgjhjkgfsdk",
	// }

	// res3 := &fooddelivery.Restaurent{
	// 	ID:       3,
	// 	Name:     "meghana",
	// 	Location: "kora",
	// 	About:    "nonveg sjdfdgjhjkgfsdk",
	// }

	// svc.AddRestarent(res1)
	// svc.AddRestarent(res2)
	// svc.AddRestarent(res3)
	// svc.ListResraurent()
	// svc.AddMenu(&fooddelivery.Menu{ID: 1,
	// 	RId:      2,
	// 	FoodMenu: "biryani",
	// 	Quantity: 10,
	// 	Price:    400})

	// svc.ViewMenu(1)
	// svc.ViewMenu(2)
	// //mid int, userid int, resid int, quantity int, price float32
	// svc.AddtoCart(1, 1, 1, 2, 800)

	ridesharingsystem.RSTest()
}

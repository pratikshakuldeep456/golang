package bookmyshow

import "fmt"

func TestBMS() {

	bmsInstance := BMSInstance()
	bmsInstance.InsertTestData()
	bmsInstance.SearchMovie("spider_man")
	bmsInstance.ViewShowTimes(1)

	u1 := &User{ID: 1, Name: "sufdg", Email: "vfkg"}

	data, err := bmsInstance.BookSeat(u1.ID, 1, []int32{2})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data.ID, data.ReservedAt, data.ShowTime)

	pgMethod := &UPIPayment{}
	data1, err1 := bmsInstance.ConfirmBooking(data.ID, pgMethod)
	if err1 != nil {
		fmt.Println(err1)
	} else {
		fmt.Println(data1)
	}

	data2, err2 := bmsInstance.ConfirmBooking(data.ID, pgMethod)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println(data2)
	}
}

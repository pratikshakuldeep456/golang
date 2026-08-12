package ridesharingsystem

import "fmt"

func RSTest() {

	instance := Instance()
	r1, _ := instance.RegisterRider("piku", "243546", &Location{4, 6})
	instance.RegisterDriver("sdfgu", "35465", &Location{12, 1}, AUTO)
	rideId, _ := instance.RequestRide(r1.ID, r1.Location, &Location{10, 4}, AUTO)
	//get driver
	findDriver := instance.FindDriver(r1.Location, AUTO)
	fmt.Println("driver", findDriver)

	instance.AcceptRide(int(rideId), findDriver.ID, AUTO)

	_, err1 := instance.AcceptRide(int(rideId), findDriver.ID, CAB)
	if err1 != nil {
		fmt.Println(err1)
	}
	instance.CalculateFare(rideId)
	// _, err = instance.AcceptRide(int(rideId), d1.ID, CAB)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	instance.CalculateFare(rideId)

}

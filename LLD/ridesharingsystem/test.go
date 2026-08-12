package ridesharingsystem

// import "fmt"

// func RSTest() {

// 	instance := Instance()
// 	r1, _ := instance.RegisterRider("piku", "243546", &Location{4, 6})
// 	instance.RegisterDriver("sdfgu", "35465", &Location{12, 1}, AUTO)
// 	rideId, _ := instance.RequestRide(r1.ID, r1.Location, &Location{10, 4}, AUTO)
// 	//get driver
// 	findDriver := instance.FindDriver(r1.Location, AUTO)
// 	fmt.Println("driver", findDriver)

// 	instance.AcceptRide(int(rideId), findDriver.ID, AUTO)

// 	_, err1 := instance.AcceptRide(int(rideId), findDriver.ID, CAB)
// 	if err1 != nil {
// 		fmt.Println(err1)
// 	}
// 	instance.CalculateFare(rideId)
// 	// _, err = instance.AcceptRide(int(rideId), d1.ID, CAB)
// 	// if err != nil {
// 	// 	fmt.Printf(err.Error())
// 	// }

// 	instance.CalculateFare(rideId)

// }

import "fmt"

func RSTest() {
	instance := Instance()

	fmt.Println("=== TEST 1: Register ===")
	r1, _ := instance.RegisterRider("piku", "243546", &Location{4, 6})
	d1, _ := instance.RegisterDriver("sdfgu", "35465", &Location{5, 5}, AUTO)
	fmt.Println("rider:", r1.ID, "driver:", d1.ID, "status:", d1.Status)

	fmt.Println("\n=== TEST 2: Request ride — EXPECT this to wrongly succeed since FindDriver isn't called ===")
	rideId, err := instance.RequestRide(r1.ID, r1.Location, &Location{10, 4}, AUTO)
	fmt.Println("rideId:", rideId, "err:", err)
	fmt.Println("BUG CHECK: was any driver actually matched/reserved? driver status:", instance.Driver[d1.ID].Status)

	fmt.Println("\n=== TEST 3: Request ride with NO eligible driver nearby — should fail but currently won't ===")
	rideId2, err := instance.RequestRide(r1.ID, r1.Location, &Location{500, 500}, AUTO)
	fmt.Println("rideId2:", rideId2, "err:", err, "(BUG: this should have failed — no driver check happens at all)")

	fmt.Println("\n=== TEST 4: FindDriver directly — will PANIC if no eligible driver exists ===")
	// deliberately search for a ride type with zero registered drivers
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("PANIC CAUGHT (expected, confirms the bug):", r)
			}
		}()
		instance.FindDriver(r1.Location, CAB)
	}()

	fmt.Println("\n=== TEST 5: CalculateFare — will silently return 0 due to wrong map check ===")
	fare := instance.CalculateFare(rideId)
	fmt.Println("fare:", fare, "(BUG CHECK: is this 0 because rID isn't in rs.Riders, even though the ride IS valid?)")

	fmt.Println("\n=== TEST 6: AcceptRide happy path ===")
	msg, err := instance.AcceptRide(int(rideId), d1.ID)
	fmt.Println("msg:", msg, "err:", err, "driver status after:", instance.Driver[d1.ID].Status)
}

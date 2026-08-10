package svc

import (
	"fmt"
)

func test() {
	lot := ParkingLotInstance()

	// 1. Add a floor
	err := lot.AddFloor(1)
	fmt.Println("AddFloor(1) err:", err)

	// 2. Add spots on floor 1
	err = lot.AddSpot(1, 101, CAR)
	fmt.Println("AddSpot(1, 101, CAR) err:", err)

	err = lot.AddSpot(1, 102, MOTORCYCLE)
	fmt.Println("AddSpot(1, 102, MOTORCYCLE) err:", err)

	// 3. Check what's actually inside lot.Floors — does floor "1" really exist as key 1?
	fmt.Printf("Floors map: %+v\n", lot.Floors)

	// 4. Try to park a car on floor 1, spot 101
	msg, err := lot.ParkVehicle(1, 101, CAR)
	fmt.Println("ParkVehicle(1, 101, CAR):", msg, err)

	// 5. Try parking a second car on the SAME spot — should fail (already occupied)
	msg, err = lot.ParkVehicle(1, 101, CAR)
	fmt.Println("ParkVehicle(1, 101, CAR) again — should fail:", msg, err)

	// 6. List available spots
	fmt.Println("Available spots:")
	lot.FindAvailableSpots()

	// 7. Grab the ticket we just created to test unpark
	// (Tickets map is keyed by ticket ID — we don't know the ID unless GenerateTicket printed/returned it)
	var createdTicket *Ticket
	for _, t := range lot.Tickets {
		createdTicket = t
		break
	}
	if createdTicket == nil {
		fmt.Println("No ticket found — ParkVehicle likely failed above.")
		return
	}
	fmt.Printf("Created ticket: %+v\n", *createdTicket)

	// 8. Unpark using that ticket
	payment := &UPIPayment{}
	msg, err = lot.UnparkVehicle(*createdTicket, payment)
	fmt.Println("UnparkVehicle:", msg, err)

	// 9. Check spot is available again
	fmt.Println("Available spots after unpark:")
	lot.FindAvailableSpots()
}

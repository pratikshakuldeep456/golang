package svc

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type ParkingLot struct {
	Floors  map[int]*ParkingFloor
	Tickets map[int]*Ticket
	// PriceSvc
	mu sync.RWMutex
}

var instance *ParkingLot

var once sync.Once

// So Singleton guarantees: one shared instance, one source of truth, accessed by every gate/goroutine in the process.
func ParkingLotInstance() *ParkingLot {
	once.Do(func() {
		instance = &ParkingLot{
			Floors:  make(map[int]*ParkingFloor),
			Tickets: make(map[int]*Ticket),
		}
	})
	return instance
}

func (pl *ParkingLot) GenerateTicket(vtype VehicleType, spot, floor int) *Ticket {

	t1 := &Ticket{
		ID:          GenerateID(),
		VehicleType: vtype,
		EntryTime:   time.Now(),
		SpotID:      spot,
		FloorID:     floor,
	}
	pl.Tickets[int(t1.ID)] = t1
	return t1

}
func (pl *ParkingLot) ParkVehicle(fid, sid int, vtype VehicleType) (string, error) {

	pl.mu.Lock()
	defer pl.mu.Unlock()
	spotStatus := pl.IsSpotAvailable(fid, sid)

	if !spotStatus {
		return "", ErrNoSpotAvailable
	}
	// status, err := pl.IsTicketValid(int(t1.ID))
	// if status == false {
	// 	return "", err
	// }

	//park

	pl.Floors[fid].Spots[sid].Status = OCCUPIED

	t1 := pl.GenerateTicket(vtype, sid, fid)
	pl.Tickets[int(t1.ID)] = t1

	return VehicleParked, nil
}

func (pl *ParkingLot) UnparkVehicle(t1 Ticket, payment Payment) (string, error) {
	status, err := pl.IsTicketValid(int(t1.ID))
	if status == false {
		return "", err
	}

	spotStatus := pl.IsSpotAvailable(t1.FloorID, t1.SpotID)

	if spotStatus {
		return "", ErrSpotAvaialbe
	}

	// 3. calculate fee
	amount := pl.CalculateFee(t1.VehicleType, t1.EntryTime, time.Now())
	// 4. initiate payment
	data := pl.InitiatePayment(&t1, amount, payment)
	if data != nil {
		return "", errors.New("failed")
	}
	// 5. if payment fails -> return error, leave spot occupied
	// 6. if payment succeeds -> free spot, update ticket.PaymentStatus = SUCCESS, ticket.ExitTime = exitTime
	// 7. return success

	pl.Floors[t1.FloorID].Spots[t1.SpotID].Status = AVAILABLE
	return SpaceMarkedAvailable, nil
}

// func (pl *ParkingLot) IsSpotAvailable(fId, sId int) bool {
// 	spot1, ok := pl.Floors[fId].Spots[sId]

// 	if !ok || spot1.Status == OCCUPIED {
// 		return true
// 	}
// 	return false

// }
func (pl *ParkingLot) IsSpotAvailable(fId, sId int) bool {
	floor, ok := pl.Floors[fId]
	if !ok {
		return false
	}
	spot, ok := floor.Spots[sId]
	if !ok || spot.Status == OCCUPIED {
		return false
	}
	return true
}
func (pl *ParkingLot) FindAvailableSpots() {

	for _, floor := range pl.Floors {
		for _, spot := range floor.Spots {
			if pl.IsSpotAvailable(int(floor.ID), int(spot.ID)) {
				fmt.Printf(" floorid: %d potid: %d\n", floor.ID, spot.ID)
			}
		}
	}
}

// func (pl *ParkingLot) Exit() {

// }

func (pl *ParkingLot) IsTicketValid(tId int) (bool, error) {
	if pl.Tickets[tId] == nil {
		return false, ErrTicketNotFound
	}
	return true, nil
}

func (pl *ParkingLot) CalculateFee(vtype VehicleType, entry, exit time.Time) int32 {
	return 100

}
func (pl *ParkingLot) InitiatePayment(ticket *Ticket, amount int32, payment Payment) error {
	paymentSvc := NewPaymentMethod(payment)
	_, err := paymentSvc.MakePayment(amount, int(ticket.ID))

	if err != nil {

		return ErrPaymentFailed
	}

	return nil
}

// admin should be able to dynamically add floors and spots.
func (pl *ParkingLot) AddFloor(fid int) error {
	if pl.Floors[fid] != nil {
		return errors.New("floor already added")
	}
	floor := &ParkingFloor{
		ID:    int32(fid),
		Spots: make(map[int]*ParkingSpot),
	}
	pl.Floors[int(floor.ID)] = floor
	return nil
}

func (pl *ParkingLot) AddSpot(fid, sid int, vtype VehicleType) error {
	if pl.Floors[fid] != nil {
		return errors.New("floor doenst added")
	}

	_, exists := pl.Floors[fid].Spots[sid]
	if exists {
		return errors.New("spots already added")
	}
	pl.Floors[fid].Spots[sid] = &ParkingSpot{ID: int32(sid),
		Status: AVAILABLE, VehicleType: vtype}

	return nil

}

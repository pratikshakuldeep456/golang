package svc

import "time"

// entities n attirbutes

type VehicleType string

const (
	CAR        VehicleType = "car"
	MOTORCYCLE VehicleType = "bike"
	TRUCK      VehicleType = "truck"
)

type ParkingStatus string

const (
	AVAILABLE ParkingStatus = "available"
	OCCUPIED  ParkingStatus = "occupied"
)

type PaymentStatus string

const (
	SUCCESS ParkingStatus = "success"
	FAILED  ParkingStatus = "failed"
)

type Ticket struct {
	ID          int32
	VehicleType VehicleType
	EntryTime   time.Time
	SpotID      int
	FloorID     int
}

// occupied    atomic.Bool  // replaces Status field
type ParkingSpot struct {
	ID          int32
	Status      ParkingStatus
	VehicleType VehicleType
}

type ParkingFloor struct {
	ID    int32
	Spots map[int]*ParkingSpot
}

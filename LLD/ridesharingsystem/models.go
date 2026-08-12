package ridesharingsystem

import "time"

type Rider struct {
	ID       int
	Name     string
	PhoneNo  string
	Location *Location
}

type Driver struct {
	ID       int
	Location *Location
	PhoneNo  string
	RideType RideType
	Status   DriverStatus
}

type Location struct {
	Lat  float64
	Long float64
}
type Ride struct {
	ID          int
	RequestedBy *Rider
	Src         *Location
	Des         *Location
	Status      RideStatus
	Fare        int32
	RideType    RideType
	AcceptedBy  *Driver
	CreatedAt   time.Time
}

type DriverStatus string

const (
	BUSY      DriverStatus = "busy"
	AVAILABLE DriverStatus = "available"
)

type RideStatus string

const (
	ONGOING   RideStatus = "ongoing"
	COMPLETED RideStatus = "completed"
	ACCEPTED  RideStatus = "accepted"
	CANCELLED RideStatus = "cancelled"
	REQUESTED RideStatus = "requested"
)

type RideType string

const (
	AUTO RideType = "auto"
	CAB  RideType = "cab"
)

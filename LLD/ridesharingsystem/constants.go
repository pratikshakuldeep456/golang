package ridesharingsystem

import "sync/atomic"

const (
	DistanceLimit     int = 10000
	NoDriveravailable     = "driver not available"
)

var riderIDCounter int32
var driverIDCounter int32
var rideIDCounter int32

func generateRiderID() int32  { return atomic.AddInt32(&riderIDCounter, 1) }
func generateDriverID() int32 { return atomic.AddInt32(&driverIDCounter, 1) }
func generateRideID() int32   { return atomic.AddInt32(&rideIDCounter, 1) }

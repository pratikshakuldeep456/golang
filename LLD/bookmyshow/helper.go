package bookmyshow

import "sync/atomic"

var id int32

func IDGenerator() int32 {
	return atomic.AddInt32(&id, 1)
}

var seatNo int32

func SeatIDGenerator() int32 {
	return atomic.AddInt32(&seatNo, 1)
}

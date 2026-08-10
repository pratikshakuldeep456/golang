package svc

import "sync/atomic"

var id int32

func GenerateID() int32 {
	return int32(atomic.AddInt32(&id, 1))
}

package cache

import (
	"errors"
	"sync/atomic"
)

var id int32

func GenerateID() int {
	return int(atomic.AddInt32(&id, 1))
}

var (
	ErrKeyNotAvailable = errors.New("key doesnt exist")
	ErrCacheEmpty      = errors.New("no data is prensent")
)

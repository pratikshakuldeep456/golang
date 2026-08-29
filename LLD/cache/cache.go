package cache

import (
	"sync"
	"time"
)

type Eviction interface {
	Get(key int) (int, error)
	Insert(key, val int) (string, error)
	Evict() (int, error)
}

type CacheElement struct {
	Val     int
	EntryAt time.Time
}
type LRU struct {
	Cap     int
	Store   map[int]*CacheElement
	Counter int
	//order   *list.List
}

var Once sync.Once
var LRUInstance *LRU

func LRUStrategy() *LRU {
	Once.Do(func() {
		LRUInstance = &LRU{
			Cap:     10,
			Store:   make(map[int]*CacheElement),
			Counter: 0,
		}
	})
	return LRUInstance
}

func (lru *LRU) Get(key int) (int, error) {
	data, exists := lru.Store[key]
	if !exists {
		return -1, ErrKeyNotAvailable
	}

	return data.Val, nil
}

func (lru *LRU) Insert(key, val int) (string, error) {
	//check if exis

	id, exists := lru.Store[key]
	if exists {
		id.EntryAt = time.Now()
		id.Val = val
		return " data added", nil
	}

	if len(lru.Store) == lru.Cap {
		//evice
		lru.Evict()
	}

	element := &CacheElement{val, time.Now()}
	lru.Store[key] = element

	return "success", nil
}

func (lru *LRU) Evict() (int, error) {
	if len(lru.Store) == 0 {
		return 0, ErrCacheEmpty
	}

	lastAccessed := lru.Store[0].EntryAt
	var id int
	for i, j := range lru.Store {
		if j.EntryAt.Before(lastAccessed) {
			id = i
		} else {
			continue
		}
	}
	delete(lru.Store, id)
	return id, nil

}

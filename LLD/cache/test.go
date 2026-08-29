package cache

import "fmt"

func CacheTest() {

	lru := LRUStrategy()
	data, err := lru.Insert(1, 1)
	if err != nil {
		fmt.Println("data", data)
	}
	fmt.Println("operation 1", data)
	lru.Insert(GenerateID(), 1)
	lru.Get(1)
	lru.Insert(5, 1)

	l := LFUStrategy()
	l.Insert(1, 10)
	l.Insert(2, 10)
	l.Get(1)

	l.Get(1)
	l.Get(1)
	l.Get(1)
	l.Insert(1, 100)
	l.Insert(3, 10)

	// l.Insert(3, 10)
	// l.Insert(4, 10)

	// l.Get(1)     // freq(1) = 2
	// l.Get(1)     // freq(1) = 3
	// l.Put(2, 20) // freq(2) = 1, the new minimum

	// l.Put(3, 30) // cache full at cap 3? no — cap is 3, only 2 keys so far, no eviction yet
}

package cache

import (
	"fmt"
	"sync"
)

// var ErrKeyNotAvailable = errors.New("key not found")
// var ErrCacheEmpty = errors.New("cache is empty")

type LFU struct {
	Cap     int
	Store   map[int]*Entry
	FreqMap map[int]map[int]struct{}
	MinFreq int
}

type Entry struct {
	val, freq int
}

var lfuOnce sync.Once
var lfuInstance *LFU

func LFUStrategy() *LFU {
	lfuOnce.Do(func() {
		lfuInstance = &LFU{
			Cap:     2,
			Store:   make(map[int]*Entry),
			FreqMap: map[int]map[int]struct{}{},
		}
	})
	return lfuInstance
}

func (l *LFU) Bump(key int, entry *Entry) {
	freq := entry.freq
	delete(l.FreqMap[freq], key)

	// fixed: check the bucket THIS key just left, not len(FreqMap)
	if len(l.FreqMap[freq]) == 0 && l.MinFreq == freq {
		l.MinFreq++
	}

	entry.freq++
	l.addToFreqSet(entry.freq, key)
}

func (l *LFU) addToFreqSet(freq, key int) {
	if l.FreqMap[freq] == nil {
		l.FreqMap[freq] = make(map[int]struct{})
	}
	l.FreqMap[freq][key] = struct{}{}
}

func (lfu *LFU) Get(key int) (int, error) {
	data, exists := lfu.Store[key]
	if !exists {
		return -1, ErrKeyNotAvailable
	}
	lfu.Bump(key, data)
	fmt.Printf("Get(%d) -> %d | store=%v\n", key, data.val, dump(lfu))
	return data.val, nil
}

func (lfu *LFU) Insert(key, val int) (string, error) {
	if data, exists := lfu.Store[key]; exists {
		data.val = val
		lfu.Bump(key, data)
		fmt.Printf("Insert(%d,%d) [update] | store=%v\n", key, val, dump(lfu))
		return "success", nil
	}

	// fixed: compare against len(Store), not len(FreqMap)
	if len(lfu.Store) >= lfu.Cap {
		victim, _ := lfu.Evict()
		fmt.Printf("evicted key=%d\n", victim)
	}

	entry := &Entry{val: val, freq: 1}
	lfu.Store[key] = entry
	lfu.addToFreqSet(1, key)
	lfu.MinFreq = 1

	fmt.Printf("Insert(%d,%d) [new] | store=%v\n", key, val, dump(lfu))
	return "success", nil
}

func (lfu *LFU) Evict() (int, error) {
	// fixed: added missing empty-cache guard
	if len(lfu.Store) == 0 {
		return 0, ErrCacheEmpty
	}

	var victim int
	for k := range lfu.FreqMap[lfu.MinFreq] {
		victim = k
		break
	}
	delete(lfu.FreqMap[lfu.MinFreq], victim)
	delete(lfu.Store, victim)
	return victim, nil
}

// dump — prints key:freq for every key currently in the cache, so you
// can see the result of each operation at a glance.
func dump(lfu *LFU) map[int]int {
	out := make(map[int]int)
	for k, e := range lfu.Store {
		out[k] = e.freq
	}
	return out
}

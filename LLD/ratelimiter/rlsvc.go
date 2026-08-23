package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

type RL interface {
	Allow() bool
}

// Token Bucket
// Leaky Bucket
// Fixed Window
// Sliding Window
// Redis-based distributed limiter

type Provider interface {
	CallAPI() error
}

type ProviderClient struct {
	Provider Provider
	RL       RL
}

func (pc *ProviderClient) Call() error {
	if !pc.RL.Allow() {
		return errors.New("not allowed")
	}
	//call provier
	return pc.Provider.CallAPI()
}

type TokenBucket struct {
	mu             sync.Mutex
	Cap            int
	Token          int
	RefillRate     int //token/sec
	LastRefillTime time.Time
}

func (tb *TokenBucket) Allow() bool {

	tb.mu.Lock()
	defer tb.mu.Unlock()

	//time last filled
	elaspedtime := time.Since(tb.LastRefillTime).Seconds()
	tb.Token += int(elaspedtime * float64(tb.RefillRate))
	if int(tb.Token) > tb.Cap {
		tb.Token = tb.Cap
	}

	tb.LastRefillTime = time.Now()

	if tb.Token < 1 {
		return false

	}

	tb.Token--
	return true

}

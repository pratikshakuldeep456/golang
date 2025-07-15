package notes

import (
	"fmt"
	"time"
)

func Limit() {
	limiter := time.Tick(500 * time.Microsecond)

	for i := 0; i < 5; i++ {
		<-limiter
		fmt.Println("limiter  for ", i, "at", time.Now())
	}
	//	limiter := time.NewTicker(1 * time.Second)
	//
	// defer limiter.Stop()
	//
	//	for i := 0; i < 5; i++ {
	//		<-limiter.C
	//		fmt.Println("limiter  for ", i, "at", time.Now())
	//	}
}

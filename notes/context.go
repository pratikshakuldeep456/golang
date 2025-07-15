package notes

import (
	"context"
	"time"
)

func Context() {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	go func() {
		// ctx.Done()
		// fmt.Println("done")
	}()
	// time.Sleep(time.Second)
	cancel()
}

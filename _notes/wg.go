package notes

import (
	"fmt"
	"sync"
	"time"
)

func WaitGrp() {

	var wg sync.WaitGroup

	tasks := []int{1, 2, 3}

	for _, j := range tasks {
		wg.Add(1)
		go func() {
			// defer wg.Done()
			time.Sleep(4 * time.Second)
			fmt.Println("processing task ", j)
		}()
	}

	wg.Wait()
	fmt.Println("All tasks completed")
}

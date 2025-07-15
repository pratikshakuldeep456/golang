package notes

//ytd Atomic Counters,Stateful Goroutines
import (
	"fmt"
	"sync"
)

func Counter() {

}
func Mutex() {
	var mu sync.Mutex
	// count := 0

	// mu.Lock()
	// go func() { count++ }()
	// mu.Unlock()

	// fmt.Println("Count:", count)
	count := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++ // Race condition!
			mu.Unlock()

		}()
	}

	wg.Wait()
	fmt.Println("Final Count:", count) // Output is unpredictable
}

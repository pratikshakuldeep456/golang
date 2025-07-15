package notes

import "fmt"

// worker pool - manages a fix no of goroutines without spawning unbounded goroutines(Creating too many goroutines without limits)

func worker(id int, job <-chan int, res chan<- int) {

	for i := range job {
		fmt.Println("processing job with id ", id)
		res <- i
	}

}
func WorkerPool() {

	job := make(chan int, 3)
	ans := make(chan int, 5)

	for i := 0; i < 3; i++ {
		go worker(i, job, ans)
	}

	for i := 0; i < 5; i++ {
		job <- i
		fmt.Println("res with i", i, " ", <-ans)
	}
	close(job)

	// for i := 0; i < 5; i++ {
	// 	fmt.Println("res with i", i, " ", <-ans)
	// }
}

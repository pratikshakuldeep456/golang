package notes

import "fmt"

func NonBlockingCh() {

	ch := make(chan int, 4)
	// ch1 := make(chan int)
	// ch2 := make(chan int)

	// go func() { ch <- 1 }()

	for j := 0; j < 3; j++ {
		ch <- j
		fmt.Println("sent job", j)
	}
	ch <- 4
	close(ch)

	for i := range ch {
		fmt.Println("i ", i)
	}
	// ch <- 4
	// ch1 <- 2
	// ch2 <- 3
	// go func() { ch <- 4 }()
	// time.Sleep(2 * time.Second)
	select {
	case msg := <-ch:
		fmt.Print("data received", msg)
	// case mag := <-ch1:
	// 	fmt.Print("data received", mag)
	// case mbg := <-ch2:
	// 	fmt.Print("data received", mbg)
	default:
		fmt.Println("no data rec")
		//
	}
}

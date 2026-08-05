package notes

import (
	"fmt"
	"time"
)

func Main2() {

	ch := make(chan string, 2)
	ch <- "110"
	ch <- "110"
	go func() {
		ch <- "4"
	}()
	go func() {
		ch <- "110"
	}()
	go func() {
		ch <- "110"
	}()
	time.Sleep(2 * time.Second)
	select {
	case msg := <-ch:
		fmt.Println("data received", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("timeout")

	}
}

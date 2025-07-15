package notes

import "fmt"

func Sendchan(ping chan<- string, m string) {
	ping <- m
}
func Recchan(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}
func Chans() {
	fmt.Println("Hello, 世界")
	ping := make(chan string, 1)
	pong := make(chan string, 1)
	Sendchan(ping, "kkd")
	Recchan(ping, pong)
	data := <-pong
	fmt.Println(data)
}

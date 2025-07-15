package main

import (
	"fmt"
)

func counter() func() int {
	count := 0
	fmt.Println("Counter initialized", count)
	return func() int {
		count++

		return count
	}

}

// func getAddress() unsafe.Pointer {
// 	x := 100
// 	return unsafe.Pointer(&x)
// }

//	func main() {
//		addr := getAddress()
//		fmt.Println(*(*int)(addr)) // Works but compiler promotes x to heap
//	}
func main() {
	c := counter()
	fmt.Println(c()) // Output: 1
	fmt.Println(c()) // Output: 2
	c1 := counter()
	fmt.Println(c1()) // Output: 1
	fmt.Println(c1()) // Output: 2
}

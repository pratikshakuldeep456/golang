package main

import (
	"fmt"
	"pratikshakuldeep456/golang/notes"
)

func Counter() func() int {
	count := 0
	fmt.Println("Counter initialized", count)
	return func() int {
		count++

		return count
	}

}

func main() {
	notes.Main2()
}

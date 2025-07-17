package main

import (
	"fmt"
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
	// notes.Main2()
	// notes.Context()
	// notes.NonBlockingCh()
	// notes.WorkerPool()
	// notes.WaitGrp()
	// notes.Limit()
	// notes.Mutex()
	// notes.FormatString()
	// notes.Json()
	// notes.Parse()
	// notes.ReadWrite()
	// notes.ContextFunc()

	// cmd := exec.Command("echo", "Hello from child process")
	// cmd.Run() // Spawns the process
}

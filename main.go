package main

import (
	"fmt"
	"runtime"
	"week-two-tasks/tasks"
)

func main() {
	tasks.Task7()
	fmt.Println("Goroutine num: ", runtime.NumGoroutine())
}

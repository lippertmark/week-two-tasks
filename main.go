package main

import (
	"fmt"
	"runtime"
	"week-two-tasks/tasks"
)

func main() {
	tasks.Task1()
	fmt.Println("Goroutine num: ", runtime.NumGoroutine())
}

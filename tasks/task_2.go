package tasks

import (
	"fmt"
	"sync"
)

func taskTwoMethodOne() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			fmt.Printf("I'm goroutine number %d(method 1)\n", num)
		}(i)
	}
	wg.Wait()
}

func taskTwoMethodTwo() {
	results := make(chan string)
	for i := 0; i < 5; i++ {
		go func(results chan string, num int) {
			results <- fmt.Sprintf("I'm goroutine number %d(method 2)", num)
		}(results, i)
	}

	for i := 0; i < 5; i++ {
		fmt.Println(<-results)
	}
}

func taskTwoMethodThreeCorrectOrder() {
	order := make([]chan struct{}, 5)
	for i := 0; i < 5; i++ {
		order[i] = make(chan struct{})
	}
	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			if num != 0 {
				<-order[num-1]
			}
			fmt.Printf("I'm goroutine number %d(method 3)\n", num)
			close(order[num])
		}(i)
	}
	wg.Wait()
}

func Task2() {
	taskTwoMethodOne()
	taskTwoMethodTwo()
	taskTwoMethodThreeCorrectOrder()
}

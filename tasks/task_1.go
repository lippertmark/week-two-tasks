package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//Задача: Напишите функцию, которая запускает горутину, выполняющую fmt.Println("Hello from goroutine!"), и использует sync.WaitGroup для ожидания её завершения.
//Какие способы есть ещё кроме waitGroup, чтобы дождаться выполнения горутины? Приведи хотя бы 2 примера.

func methodOne() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("Hello from goroutine! method 1")
	}(wg)
	wg.Wait()
}

func methodTwo() {
	ch := make(chan struct{})
	go func(ch chan<- struct{}) {
		fmt.Println("Hello from goroutine! method 2")
		ch <- struct{}{}
	}(ch)
	_ = <-ch
}

func methodThree() {
	go func() {
		fmt.Println("Hello from goroutine! method 3")
	}()
	time.Sleep(time.Second)
}

func methodFour() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("Hello from goroutine! method 4")
		cancel()
	}()
	select {
	case <-ctx.Done():
	}

}

func methodFive() {
	go func() {
		fmt.Println("Hello from goroutine! method 5")
	}()
	fmt.Scanln()
}

func Task1() {
	methodOne()
	methodTwo()
	methodThree()
	methodFour()
	methodFive()

}

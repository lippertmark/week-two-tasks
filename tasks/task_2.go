package tasks

import (
	"fmt"
	"sync"
)

//Задача: Напишите программу, которая запускает 5 горутин, каждая из которых печатает свой номер (от 1 до 5), и использует sync.WaitGroup для их синхронизации(нужно подождать их выполнения).
//
//Можно ли решить задачу без waitGroup? Какие есть варианты?
//Можно ли сделать так чтобы номера выводились в определённом порядке? Почему? Может всё-таки можно?
//
//Как влияет GOMAXPROCS на выполнение программы?

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

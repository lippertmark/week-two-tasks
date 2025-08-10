package tasks

import (
	"fmt"
	"sync"
	"time"
)

func Merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(cs))
	for _, ch := range cs {
		go func() {
			for val := range ch {
				out <- val
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// создаёт канал, который посылает числа с задержкой
func generator(from, count int, delay time.Duration) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < count; i++ {
			ch <- from + i
			time.Sleep(delay)
		}
	}()
	return ch
}

func Task9() {
	a := generator(100, 3, 200*time.Millisecond) // 100, 101, 102
	b := generator(200, 2, 300*time.Millisecond) // 200, 201

	merged := Merge(a, b)

	for val := range merged {
		fmt.Println("получено:", val)
	}
}

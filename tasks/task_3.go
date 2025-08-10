package tasks

import "fmt"

func genNums() chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 1; i <= 5; i++ {
			ch <- i
		}
	}()
	return ch
}

func task3MethodOne() {
	out := genNums()
	sum := 0
	for val := range out {
		sum += val
	}
	fmt.Println(sum)
}

func task3MethodTwo() {
	out := genNums()
	sum := 0
	for {
		select {
		case val, ok := <-out:
			if !ok {
				fmt.Println(sum)
				return
			}

			sum += val
		}
	}

}
func Task3() {
	task3MethodOne()
	task3MethodTwo()
}

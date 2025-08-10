package tasks

import "fmt"

func Split(in <-chan int, n int) []<-chan int {
	result := make([]<-chan int, n)
	for id := range n {
		ch := make(chan int)
		go func(id int) {
			defer close(ch)
			for val := range in {
				ch <- val
			}
		}(id)
		result[id] = ch
	}
	return result
}

func Task10() {
	in := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			in <- i
		}
		close(in)
	}()

	outChans := Split(in, 4)

	done := make(chan bool)
	for i, ch := range outChans {
		go func(id int, ch <-chan int) {
			for val := range ch {
				fmt.Printf("Channel #%d got: %d\n", id, val)
			}
			done <- true
		}(i, ch)
	}

	for i := 0; i < len(outChans); i++ {
		<-done
	}
	close(done)

}

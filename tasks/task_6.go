package tasks

import (
	"context"
	"fmt"
	"time"
)

func StartBatchProcessor(ctx context.Context, input <-chan int) {
	var batch []int
	timeout := 2 * time.Second
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Canceled")
			return
		case <-timer.C:
			fmt.Println("timer")
			fmt.Println("Processed batch:", batch)
			batch = nil
			timer.Reset(timeout)
		case val, ok := <-input:
			if !ok {
				return
			}
			batch = append(batch, val)
			if len(batch) == 5 {
				fmt.Println("Processed batch:", batch)
				batch = nil
			}
			timer.Stop()
			timer = time.NewTimer(timeout)

		}
	}

}

func Task6() {
	input := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go StartBatchProcessor(ctx, input)

	go func() {
		for i := 1; i <= 20; i++ {
			input <- i
			time.Sleep(300 * time.Millisecond)
		}
	}()

	<-ctx.Done()
	fmt.Println("Main: processing stopped")

	//input := make(chan int)
	//ctx, cancel := context.WithCancel(context.Background())
	//go func() {
	//	for i := 1; i <= 20; i++ {
	//		input <- i
	//		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	//	}
	//	close(input)
	//}()
	//
	//go StartBatchProcessor(ctx, input)
	//
	//time.Sleep(1 * time.Second)
	//time.Sleep(15 * time.Second)
	//cancel()
	//fmt.Println("cancel")
	//
	//time.Sleep(1 * time.Second)
	//fmt.Println("finish")
}

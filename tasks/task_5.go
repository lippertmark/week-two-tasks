package tasks

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//Задача: Напишите программу, где 10 горутин инкрементируют один счётчик без использования мютексов, через атомики.
//
//Что, если не использовать атомик? Что лучше, атомик или мютекс?

func inc(wg *sync.WaitGroup, cnt *int32) {
	defer wg.Done()
	atomic.AddInt32(cnt, 1)
}

func Task5() {
	var cnt int32
	wg := &sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go inc(wg, &cnt)
	}
	wg.Wait()
	fmt.Println(cnt)
}

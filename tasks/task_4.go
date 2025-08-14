package tasks

import (
	"fmt"
	"sync"
)

//Задача: Напишите программу, где 10 горутин инкрементируют один счётчик, защищая его sync.Mutex.
//
//1. Что если не обложить мютексом? Воспроизвести race condition.

func Task4() {
	mtx := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	cnt := 0
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(mtx *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			mtx.Lock()
			cnt += 1
			mtx.Unlock()
		}(mtx, wg)
	}
	wg.Wait()
	fmt.Println(cnt)
}

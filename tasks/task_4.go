package tasks

import (
	"fmt"
	"sync"
)

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

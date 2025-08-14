package tasks

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// задачка из собеседований https://www.notion.so/Golang-2118033d6fe580c4ad3ef7c161593172?source=copy_link#2118033d6fe581a69a6ee0c5040c7fb1

func Task11() {
	//users := make([]User, 1000)
	//for i := 0; i < 1000; i++ {
	//	users = append(users, User{strconv.Itoa(i)})
	//
	users := []User{{"aaa"}, {"bbb"}, {"ccc"}, {"ddd"}, {"eee"}}
	resp, err := Do(context.Background(), users)
	fmt.Println(resp, err)
	fmt.Println(len(resp))
}

type User struct {
	Name string
}

func fetchByName(ctx context.Context, userName string) (int, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	select {
	case <-ctxWithTimeout.Done():
		//if userName == "bbb" {
		//	return 0, errors.New("some error(bbb)")
		//}
		//if userName == "ccc" {
		//	return 0, errors.New("some error(ccc)")
		//}
		return rand.Int() % 100000, nil
	case <-ctx.Done():
		return 0, errors.New("operation canceled")

	}

}

func Do(ctx context.Context, users []User) (map[string]int, error) {
	collected := make(map[string]int)
	errorCh := make(chan error)
	mu := &sync.Mutex{}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg := &sync.WaitGroup{}
	for _, u := range users {
		wg.Add(1)
		go func(user User) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			userID, err := fetchByName(ctx, user.Name)
			if err != nil {
				select {
				case errorCh <- err:
					cancel()
				case <-ctx.Done():
				}
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			mu.Lock()
			collected[user.Name] = userID
			mu.Unlock()
		}(u)
	}
	go func() {
		wg.Wait()
		cancel()
	}()

	select {
	case <-ctx.Done():
		return collected, nil
	case err := <-errorCh:
		return map[string]int{}, err
	}
}

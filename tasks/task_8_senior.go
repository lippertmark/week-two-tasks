package tasks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// 4. Ограничить общее количество параллельных запросов (например, максимум 10 одновременно через семафор — chan struct{}).
// А лучше использовать для этого паттерн worker pool Go by Example: Worker Pools

const maxBodyLengthSenior = 100
const timeoutOfRequestSenior = 1 * time.Second

func writeResponseSenior(mtx *sync.Mutex, url, response string, results map[string]string) {
	mtx.Lock()
	results[url] = response
	mtx.Unlock()
}

func fetchURLSenior(mtx *sync.Mutex, ctx context.Context, client *http.Client, url string, results map[string]string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeoutOfRequestSenior)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxTimeout, "GET", url, nil)
	if err != nil {
		writeResponseSenior(mtx, url, err.Error(), results)
		return errors.New("error 1")
	}
	resp, err := client.Do(req)
	if err != nil {
		writeResponseSenior(mtx, url, err.Error(), results)
		return errors.New("error 2")
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("ERROR")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeResponseSenior(mtx, url, err.Error(), results)
		return errors.New("error 3")
	}
	trimmed := body
	if len(body) > maxBodyLengthSenior {
		trimmed = body[:maxBodyLengthSenior]
	}
	writeResponseSenior(mtx, url, fmt.Sprintf("code: %v: body: %s", resp.StatusCode, trimmed), results)
	return nil
}

func FetchURLsSenior(urls []string, parallelism int) map[string]string {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}
	results := make(map[string]string)
	startTime := time.Now()
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}
	defer transport.CloseIdleConnections()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workers := make(chan int, parallelism)
	for workerId := 0; workerId < parallelism; workerId++ {
		workers <- workerId
	}
	for _, url := range urls {
		select {
		case workerId := <-workers:
			wg.Add(1)
			go func(u string, workerId int) {
				defer wg.Done()
				err := fetchURLSenior(mtx, ctx, client, u, results)
				if err != nil {
					cancel()
				}
				workers <- workerId
			}(url, workerId)
		case <-ctx.Done():
			break
		}

	}
	wg.Wait()
	fmt.Println("Took", time.Since(startTime))
	return results
}

func Task8Senior() {
	urls := []string{
		"https://example.com",
		"https://httpbin.org/get",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404",
		"https://httpbin.org/delay/2",
		"https://golang.org",
		"https://api.github.com",
		"https://httpbin.org/bytes/50",
		"https://httpbin.org/status/500",
		//"https://invalid.url",
	}
	result := FetchURLsSenior(urls, 1)
	for key, value := range result {
		fmt.Printf("%v %v\n------------------------\n", key, value)
	}
	fmt.Println(len(result))
}

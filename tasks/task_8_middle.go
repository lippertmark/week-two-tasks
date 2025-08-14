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

// 2. Переписать код, чтобы использовать контекст отмены
// 3. При первой ошибке отменять все остальные запросы

const maxBodyLengthMiddle = 100
const timeoutOfRequestMiddle = 2 * time.Second

func writeResponseMiddle(mtx *sync.Mutex, url, response string, results map[string]string) {
	mtx.Lock()
	results[url] = response
	mtx.Unlock()
}

func fetchURLMiddle(mtx *sync.Mutex, ctx context.Context, client *http.Client, url string, results map[string]string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeoutOfRequestMiddle)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxTimeout, "GET", url, nil)
	if err != nil {
		writeResponseMiddle(mtx, url, err.Error(), results)
		return errors.New("error 1")
	}
	resp, err := client.Do(req)
	if err != nil {
		writeResponseMiddle(mtx, url, err.Error(), results)
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
		writeResponseMiddle(mtx, url, err.Error(), results)
		return errors.New("error 3")
	}
	trimmed := body
	if len(body) > maxBodyLengthMiddle {
		trimmed = body[:maxBodyLengthMiddle]
	}
	writeResponseMiddle(mtx, url, fmt.Sprintf("code: %v: body: %s", resp.StatusCode, trimmed), results)
	return nil
}

func FetchURLsMiddle(urls []string) map[string]string {
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
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			err := fetchURLMiddle(mtx, ctx, client, u, results)
			if err != nil {
				cancel()
			}
		}(url)

	}
	wg.Wait()
	fmt.Println("Took", time.Since(startTime))
	return results
}

func Task8Middle() {
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
		"https://invalid.url",
	}
	result := FetchURLsMiddle(urls)
	for key, value := range result {
		fmt.Printf("%v %v\n------------------------\n", key, value)
	}
	fmt.Println(len(result))
}

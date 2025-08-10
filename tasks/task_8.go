package tasks

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const maxBodyLength = 100
const timeoutOfRequest = 7 * time.Second

func writeResponse(mtx *sync.Mutex, url, response string, results map[string]string) {
	mtx.Lock()
	results[url] = response
	mtx.Unlock()
}

func fetchURL(wg *sync.WaitGroup, mtx *sync.Mutex, client *http.Client, url string, results map[string]string) {
	defer wg.Done()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		writeResponse(mtx, url, "error 1", results)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		writeResponse(mtx, url, err.Error(), results)
		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("ERROR")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		writeResponse(mtx, url, err.Error(), results)
		return
	}
	trimmed := body
	if len(body) > maxBodyLength {
		trimmed = body[:maxBodyLength]
	}
	writeResponse(mtx, url, fmt.Sprintf("code: %v: body: %s", resp.StatusCode, trimmed), results)
}

func FetchURLs(urls []string) map[string]string {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}
	results := make(map[string]string)
	startTime := time.Now()
	transport := &http.Transport{}
	client := &http.Client{
		Timeout:   timeoutOfRequest,
		Transport: transport,
	}
	defer transport.CloseIdleConnections()
	for _, url := range urls {
		wg.Add(1)
		go fetchURL(wg, mtx, client, url, results)
	}
	wg.Wait()
	fmt.Println("Took", time.Since(startTime))
	return results
}

func Task8() {
	urls := []string{
		"https://httpbin.org/get",
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404",
		"https://httpbin.org/delay/2",
		"https://example.com",
		"https://golang.org",
		"https://api.github.com",
		"https://httpbin.org/bytes/50",
		"https://httpbin.org/status/500",
		"https://invalid.url",
	}
	result := FetchURLs(urls)
	for key, value := range result {
		fmt.Printf("%v %v\n------------------------\n", key, value)
	}
	fmt.Println(len(result))
}

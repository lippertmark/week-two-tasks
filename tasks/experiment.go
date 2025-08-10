package tasks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func Experiment() {
	transport := &http.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second, // общий таймаут запроса
	}

	resp, err := client.Get("https://www.youtube.com/watch?v=yZrGRU4mmBE")
	if err != nil {
		log.Fatalf("Ошибка запроса: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Ошибка закрытия тела ответа: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения тела: %v", err)
	}

	fmt.Println(string(body))

	transport.CloseIdleConnections()
}

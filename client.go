package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func sendRequests() {
	requestsSent := 1
	for requestsSent <= 50 {
		requestsSent += 1
		time.Sleep(50 * time.Millisecond)
		_, err := http.Get("http://localhost:8091/rate-limiter")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func client() {
	sendRequests()
	fmt.Println("Client has finished sending its request")
}

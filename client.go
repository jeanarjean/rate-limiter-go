package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func sendRequests() {
	requestsSent := 1
	for requestsSent <= 50 {
		requestsSent += 1
		time.Sleep(50 * time.Millisecond)
		response, err := http.Get("http://localhost:8091/rate-limiter")
		if err != nil {
			log.Fatalln(err)
		}
		if response.StatusCode != 200 {
			fmt.Printf("Client's request number %v has been dropped \n", requestsSent)
		}
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Program is done executing and is exiting")
	os.Exit(0)
}

func client() {
	sendRequests()
	fmt.Println("Client has finished sending its request")
}

package main

import (
	"fmt"
	"net/http"
	"time"
)

type Message struct {
	Url    string
	Header http.Header
	Method string
}

// Every 1000 seconds run the cleanup
const arraySize = 10

var requestsQueue = make(chan Message, 300)

var index = 0
var httpClient = http.Client{}
var limiter = time.Tick(500 * time.Millisecond)
var requestsReceived = 0

func throttle() {
	// I MUST HAVE BEEN SENDING AN EMPTY REQUEST HERE, FIX THIS
	for {
		<-limiter
		requestToSend := <-requestsQueue
		sendReceivedRequestToServer(requestToSend)
	}
}

func receiveHttpCall(w http.ResponseWriter, req *http.Request) {
	requestsReceived += 1
	fmt.Println("Rate Limiter has received requests", requestsReceived)

	message := Message{
		req.URL.String(),
		req.Header,
		req.Method,
	}
	requestsQueue <- message
}
func sendReceivedRequestToServer(message Message) {
	request, err := http.NewRequest(message.Method, "http://localhost:8090/server", nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header = message.Header
	_, err = httpClient.Do(request)
	if err != nil {
		fmt.Println(err)
	}
}

func rateLimiter() {
	http.HandleFunc("/rate-limiter", receiveHttpCall)

	go throttle()

	http.ListenAndServe(":8091", nil)
}

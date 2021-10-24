package main

import (
	"fmt"
	"net/http"
	"time"
)

type Message struct {
	Url       string
	Header    http.Header
	Method    string
	Timestamp time.Time
}

type Algorithm string

const (
	NoRateLimiting     = "NoRateLimiting"
	ConstantRate       = "ConstantRate"
	TokenBucket        = "TokenBucket"
	LeakingBucket      = "LeakingBucket"
	FixedWindowCounter = "FixedWindowCounter"
	SlidingWindowLog   = "SlidingWindow"
)

// Every 1000 seconds run the cleanup
const arraySize = 10

var requestsQueue = make(chan Message, 300)

var index = 0
var httpClient = http.Client{}
var requestsReceived = 0

//TODO: move this in noConstantRateLimiting scope
var limiter = time.Tick(500 * time.Millisecond)

func noRateLimiting() {
	for {
		requestToSend := <-requestsQueue
		sendReceivedRequestToServer(requestToSend)
	}
}

func constantRate() {
	for {
		<-limiter
		requestToSend := <-requestsQueue
		sendReceivedRequestToServer(requestToSend)
	}
}

func tokenBucket() {
	tokens := make(chan bool, 4)
	interval := time.Tick(500 * time.Millisecond)
	go func() {
		for {
			<-interval
			if len(tokens) != cap(tokens) {
				fmt.Printf("Adding a token in the bucket\n")
				tokens <- true
			}
			if len(tokens) != cap(tokens) {
				fmt.Printf("Adding a token in the bucket\n")
				tokens <- true
			}
		}
	}()
	for {
		<-tokens
		requestToSend := <-requestsQueue
		sendReceivedRequestToServer(requestToSend)
	}
}

func leakingBucket() {
	for {
		<-limiter
		requestToSend := <-requestsQueue
		sendReceivedRequestToServer(requestToSend)
		requestToSend = <-requestsQueue
		sendReceivedRequestToServer(requestToSend)
	}
}

func fixedWindowCounter() {
	requestsDuringInterval := 0
	const maxRequests = 4
	interval := time.Tick(1000 * time.Millisecond)
	go func() {
		for {
			<-interval
			requestsDuringInterval = 0
		}
	}()
	for {
		if requestsDuringInterval < 4 {
			requestsDuringInterval++
			requestToSend := <-requestsQueue
			sendReceivedRequestToServer(requestToSend)
		}
	}
}

func slidingWindowLog() {
	var recentRequests []Message
	for {
		recentRequest := <-requestsQueue
		var temp []Message
		for _, r := range recentRequests {
			//fmt.Printf("Request Timestamp: %v, Now: %v \n", r.Timestamp.Add((time.Second)), time.Now())
			fmt.Printf("%v \n", len(recentRequests))
			if r.Timestamp.Add(time.Second).After(time.Now()) {
				temp = append(temp, r)
			}
		}

		recentRequests = temp

		// There is ambiguity here whether dropped requests should be kept in the log, however in this program's case/benchmark, it DOS the server
		// to keep them in the log, so we won't
		if len(recentRequests) < 3 {
			recentRequests = append(recentRequests, recentRequest)
			sendReceivedRequestToServer(recentRequest)
		}
	}
}

func receiveHttpCall(w http.ResponseWriter, req *http.Request) {
	requestsReceived += 1
	fmt.Printf("Rate Limiter has %d received requests\n", requestsReceived)

	message := Message{
		req.URL.String(),
		req.Header,
		req.Method,
		time.Now(),
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

func rateLimiter(a Algorithm) {
	http.HandleFunc("/rate-limiter", receiveHttpCall)

	switch a {
	case NoRateLimiting:
		go noRateLimiting()
	case ConstantRate:
		go constantRate()
	case TokenBucket:
		go tokenBucket()
	case LeakingBucket:
		go leakingBucket()
	case FixedWindowCounter:
		go fixedWindowCounter()
	case SlidingWindowLog:
		go slidingWindowLog()
	}

	http.ListenAndServe(":8091", nil)
}

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

type Algorithm string

const (
	NoRateLimiting     = "NoRateLimiting"
	TokenBucket        = "TokenBucket"
	LeakingBucket      = "LeakingBucket"
	FixedWindowCounter = "FixedWindowCounter"
	SlidingWindowLog   = "SlidingWindow"
)

var httpClient = http.Client{}

var rateLimitingFunc func() bool = noRateLimiting()

func noRateLimiting() func() bool {
	return func() bool {
		return true
	}
}

func tokenBucket() func() bool {
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

	return func() bool {
		select {
		case <-tokens:
			return true
		default:
			return false
		}
	}
}

func leakingBucket() func() bool {
	const interval = 500 * time.Millisecond
	var lastRequestTime = time.Now().Add(-interval)

	return func() bool {
		if lastRequestTime.Add(interval).Before(time.Now()) {
			lastRequestTime = time.Now()
			return true
		}

		return false
	}
}

func fixedWindowCounter() func() bool {
	requestsDuringInterval := 0
	const maxRequests = 4
	interval := time.Tick(1000 * time.Millisecond)
	go func() {
		for {
			<-interval
			requestsDuringInterval = 0
		}
	}()

	return func() bool {
		if requestsDuringInterval < 4 {
			requestsDuringInterval++
			return true
		}

		return false
	}
}

func slidingWindowLog() func() bool {
	var recentRequests []time.Time
	const interval = 500 * time.Millisecond

	return func() bool {
		var temp []time.Time
		for _, r := range recentRequests {
			if r.Add(interval).After(time.Now()) {
				temp = append(temp, r)
			}
		}

		recentRequests = temp

		// There is ambiguity here whether dropped requests should be kept in the log, however in this program's case/benchmark, it DOS the server
		// to keep them in the log, so we won't
		if len(recentRequests) < 3 {
			recentRequests = append(recentRequests, time.Now())
			return true
		}
		return false
	}
}

func receiveHttpCall(w http.ResponseWriter, req *http.Request) {
	message := Message{
		req.URL.String(),
		req.Header,
		req.Method,
	}

	if rateLimitingFunc() {
		sendReceivedRequestToServer(message)
	} else {
		http.Error(w, "Request has been throttled", http.StatusTooManyRequests)
	}
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

	switch a {
	case NoRateLimiting:
		rateLimitingFunc = noRateLimiting()
	case TokenBucket:
		rateLimitingFunc = tokenBucket()
	case LeakingBucket:
		rateLimitingFunc = leakingBucket()
	case FixedWindowCounter:
		rateLimitingFunc = fixedWindowCounter()
	case SlidingWindowLog:
		rateLimitingFunc = slidingWindowLog()
	}

	http.HandleFunc("/rate-limiter", receiveHttpCall)

	http.ListenAndServe(":8091", nil)
}

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	SuccessColor = "\033[1;32m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

func sendRequests() {
	requestsSent := 0
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for requestsSent <= 15 {
		requestsSent += 1
		time.Sleep(time.Duration(r1.Intn(5)) * 100 * time.Millisecond)
		fmt.Printf(InfoColor, fmt.Sprintf("Sending request number %v... ", requestsSent))
		response, err := http.Get("http://localhost:8091/rate-limiter")
		if err != nil {
			log.Fatalln(err)
		}
		if response.StatusCode == 200 {
			fmt.Printf(SuccessColor, "the request has been accepted \n")
		} else {
			fmt.Printf(WarningColor, "the request has been refused \n")
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

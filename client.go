package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func hola() {
	sum := 1
	for sum < 50 {
		sum += 1
		time.Sleep(50 * time.Millisecond)
		_, err := http.Get("http://localhost:8091/rate-limiter")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func client() {
	hola()
	fmt.Println("Client has finished sending its request")
}

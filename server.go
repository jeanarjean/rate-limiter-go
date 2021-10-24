package main

import (
	"fmt"
	"net/http"
)

var serverRequestsReceived = 0

func babyWeKnow(w http.ResponseWriter, req *http.Request) {
	serverRequestsReceived += 1
	fmt.Println("Server has received %d requests", serverRequestsReceived)
}

func server() {
	http.HandleFunc("/server", babyWeKnow)

	http.ListenAndServe(":8090", nil)
}

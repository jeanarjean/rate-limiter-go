package main

import (
	"net/http"
)

var serverRequestsReceived = 0

func receiveRequest(w http.ResponseWriter, req *http.Request) {
	serverRequestsReceived++
}

func server() {
	http.HandleFunc("/server", receiveRequest)

	http.ListenAndServe(":8090", nil)
}

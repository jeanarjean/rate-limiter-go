package main

import (
	"fmt"
	"net/http"
	"os"
)

var serverRequestsReceived = 0

func receiveRequest(w http.ResponseWriter, req *http.Request) {
	serverRequestsReceived += 1
	fmt.Printf("Server has received %v requests\n", serverRequestsReceived)

	//TOOD: since some requests now get dropped, we'll need to change the logic here.
	if serverRequestsReceived > 49 {
		fmt.Println("Program is done executing and is exiting")
		os.Exit(0)
	}
}

func server() {
	http.HandleFunc("/server", receiveRequest)

	http.ListenAndServe(":8090", nil)
}

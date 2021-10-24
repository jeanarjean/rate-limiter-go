package main

func main() {
	//TOOD: Add CLI to decide which algorithms will run
	go client()
	go rateLimiter(SlidingWindowLog)
	server()
}

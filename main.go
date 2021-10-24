package main

func main() {
	// Add requests to send parameter?
	go client()
	go rateLimiter(SlidingWindowLog)
	server()
}

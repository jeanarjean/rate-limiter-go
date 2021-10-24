package main

func main() {
	go client()
	go rateLimiter()
	server()
}

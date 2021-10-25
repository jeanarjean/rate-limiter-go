package main

import (
	"fmt"
	"os"
)

func main() {
	go client()

	if len(os.Args) < 2 {
		fmt.Println("expected 'noRateLimiting', 'tokenBucket', 'leakingBucket, 'fixedWindowWCounter' or 'slidingWindowLog' to be specified")
		fmt.Println("example: '.\\rate-limiter tokenBucket'")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "noRateLimiting":
		go rateLimiter(NoRateLimiting)
	case "tokenBucket":
		go rateLimiter(TokenBucket)
	case "leakingBucket":
		go rateLimiter(LeakingBucket)
	case "fixedWindowCounter":
		go rateLimiter(FixedWindowCounter)
	case "slidingWindowLog":
		go rateLimiter(SlidingWindowLog)
	default:
		fmt.Println("expected 'noRate', 'tokenBucket', 'leakingBucket, 'fixedWindowWCounter' or 'slidingWindowLog' to be specified")
		fmt.Println("example: '.\\rate-limiter tokenBucket'")
		os.Exit(1)
	}

	server()
}

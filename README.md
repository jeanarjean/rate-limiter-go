# rate-limiter-go

A very simple rate limiter in go, made as a learning project to learn go and rate limiting patterns!

## Demo:

![Demo](/static/demo.gif?raw=true "Demo")

## Running the project:

To execute the binary, specify an algorithm:

```bash
.\rate-limiter-go.exe noRateLimiting
.\rate-limiter-go.exe tokenBucket
.\rate-limiter-go.exe leakingBucket
.\rate-limiter-go.exe fixedWindowCounter
.\rate-limiter-go.exe slidingWindowLog
```

To run from source:

```bash
go run . noRateLimiting
go run . tokenBucket
go run . leakingBucket
go run . fixedWindowCounter
go run . slidingWindowLog
```

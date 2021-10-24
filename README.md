# rate-limiter-go

A very simple rate limiter in go, made as a learning project to learn go and rate limiting patterns!

To Run from source:

```bash
go run . noRateLimiting
```

## Available algorithms

### No Rate Limiting

```sh
go run . noRateLimiting
```

### Token Bucket

```sh
go run . tokenBucket
```

### Leaking Bucket

```sh
go run . leakingBucket
```

### FixedWindowCounter

```sh
go run . fixedWindowCounter
```

### SlidingWindowLog

```sh
go run . slidingWindowLog
```

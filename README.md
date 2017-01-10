# trabago
Go library for various work-related patterns (rate-limiting, worker-pools, scheduling, etc.). 
All provided data structures are concurrency-safe.

> Work in progress (1/10/2017)

## Features

- Rate limiting
- Worker pools
- Stacks and priority queues
- ...
- No third-party package dependencies

## Requirements

- Go 1.7+

## Installation

`go get github.com/patrickrand/trabago`

## Usages

### Rate limiting

```go
func main() {
    // number of permitted executions per time-interval
    var n uint = 100

    // time-interval to rate-limit over
    duration := 10 * time.Second

    viewLimiter := trabago.NewRateLimiter(n, duration)

    http.Handle("/view", func(w http.ResponseWriter, r *http.Request) {
        if viewLimiter.IsRateLimited() {
            w.WriteHeader(http.StatusTooManyRequests)
            w.Write([]byte(`Client has exceeded API rate limit.`))
            return
        }

        ...
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`Some view...`))
    })

    http.ListenAndServe(":8080", nil)
}
```
## Rate Limiter

![Rate Limiter](https://github.com/althk/ratelimiter/actions/workflows/ci-main.yml/badge.svg)

Provides the following rate limiters:

- [Token Bucket](https://en.wikipedia.org/wiki/Token_bucket)
- [Sliding Window Counter](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/)

The package also provides the following:

### storage drivers

- Local (in-memory)
  - A rate limiter configured with this store simply stores all data locally in memory.
    It evicts all keys that have not been touched in the last minute.
- Redis backed
  - This provides a more reliable storage for the rate limiter.
  - Every key has a TTL of one minute.
  - Currently does not have TLS support.

### middleware

- http
  - a pluggable http middleware to make it easy for using the rate limiters.
  - see [example/](example) directory for usage.

### configuration

The rate limiters can be configured via the `LimiterOptions` type.
See [example/main.go](example/main.go) for hints.


### TODOs

- add TLS support for Redis store
- make the error message configurable (currently it returns a 429 status code with "Quota Exceeded" as the message)
- add GRPC middleware
- allow fetching limits per IP from a config source and use custom limits per client,
  and fall back to the global default limit if a client does not have a custom rate limit

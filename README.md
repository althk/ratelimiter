## Rate Limiter

Provides the following rate limiters:

- [Token Bucket](https://en.wikipedia.org/wiki/Token_bucket)
  - Currently, the default limit is set to 1 request per second per IP.
- [Sliding Window Counter](https://blog.cloudflare.com/counting-things-a-lot-of-different-things/)
  - Currently, the default limit is set to 1 request per second per IP.

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

### TODOs

- make the rate limiters configurable (rates, limits, etc.)
- add TLS support for Redis store.
- make the error message configurable (currently it returns a 429 status code with "Quota Exceeded" as the message)
- add GRPC middleware
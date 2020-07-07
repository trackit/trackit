# TrackIt
## Logging

Logging in TrackIt is not done through `fmt` or `log` packages. It is done using the `github.com/trackit/jsonlog` package.

The logger exposes multiple functions for logging with different logging levels:
- Debug
- Info
- Warning
- Error

Example:
```go
logger := jsonlog.LoggerFromContextOrDefault(ctx)	
logger.Info("An example message.", map[string] interface{} {
    "reportDate": time.Now().UTC(),
    "moreInfo": "Another example message.",
})
```
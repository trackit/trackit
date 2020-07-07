# Logging

Logging in TrackIt is not done through `fmt` or `log` packages. It is done using the `github.com/trackit/jsonlog` package.

The logger exposes multiple functions for logging with different logging levels:
- Debug
- Info
- Warning
- Error

Each of these functions takes two parameters: a string message, and an interface{} object to log additional data.

### Example:
```go
logger := jsonlog.LoggerFromContextOrDefault(ctx)	
logger.Info("An example message.", map[string] interface{} {
    "reportDate": time.Now().UTC(),
    "moreInfo": "Another example message.",
})
```
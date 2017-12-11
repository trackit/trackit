package routes

// contextKey represents a key in a context. Using an unexported type in this
// fashion ensures there can be no collision with a key from some other
// package.
type contextKey int

const (
	// contextKeyRequestId is the key for a request's random ID stored in
	// its context.
	contextKeyRequestId = contextKey(iota)
)

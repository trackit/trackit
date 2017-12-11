package routes

// argumentKey represents a key in an Arguments map.
type argumentKey int

const (
	// argumentKeyBody associates with a request body as produced by the
	// RequestBody decorator.
	argumentKeyBody = argumentKey(iota)
)

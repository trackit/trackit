package routes

// argumentKey represents a key in an Arguments map.
type argumentKey int

const (
	// argumentKeyJsonBody associates with a request body as produced by
	// the JsonRequestBody decorator.
	argumentKeyJsonBody = argumentKey(iota)
	// argumentKeyJsonBodyType associates with the type of a request body
	// as produced by the JsonRequestBody decorator.
	argumentKeyJsonBodyType
)

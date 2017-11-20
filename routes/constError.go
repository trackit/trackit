package routes

// constError is a string which fullfils the error interface.
type constError string

func (ce constError) Error() string {
	return string(ce)
}

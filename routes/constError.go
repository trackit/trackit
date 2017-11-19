package routes

type constError string

func (ce constError) Error() string {
	return string(ce)
}

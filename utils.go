package main

type ErrorMessage struct {
	s string
}

func (e *ErrorMessage) Error() string {
	return e.s
}

package elevatorsystem

import "errors"

var (
	InValidRequest         = errors.New("request is invalid")
	ErrNoElevatorAvailable = errors.New("invlaid req")
	ErrInvalidElevator     = errors.New("invlaid elelvator")
)

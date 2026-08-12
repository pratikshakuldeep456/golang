package ridesharingsystem

import "errors"

var (
	ErrRideNotPresentInSystem       = errors.New("ride doesnt exist")
	ErrDriverNotPresentInSystem     = errors.New("driver doesnt exist")
	ErrNoDriverAvailable            = errors.New("not avaialble")
	ErrRiderNotPresentInSystem      = errors.New("not present")
	ErrRiderBusyOrDifferentRideType = errors.New("driver is not avaialbr or ride type is differnet")
)

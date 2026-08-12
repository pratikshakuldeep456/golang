package ridesharingsystem

import "errors"

var (
	ErrRiderNotPresentInSystem      = errors.New("rider doesnt exist")
	ErrDriverNotPresentInSystem     = errors.New("driver doesnt exist")
	ErrRiderBusyOrDifferentRideType = errors.New("driver is not avaialbr or ride type is differnet")
)

package svc

import "errors"

var (
	ErrNoSpotAvailable = errors.New("no spot available for vehicle type")
	ErrTicketNotFound  = errors.New("ticket not found")
	ErrPaymentFailed   = errors.New("payment failed")
	ErrSpotAvaialbe    = errors.New("spot isnt occupied")
)

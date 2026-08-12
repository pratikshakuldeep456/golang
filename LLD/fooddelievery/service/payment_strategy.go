package service

type PaymentGateway interface {
	Pay(amount float32, orderID int) (string, error)
}

type ThirdPartyPaymentGateway struct{} // stub — assume Stripe/Razorpay/etc in reality
func (t *ThirdPartyPaymentGateway) Pay(amount float32, orderID int) (string, error) {
	return "payment successful", nil // stubbed
}

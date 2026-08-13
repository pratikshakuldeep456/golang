package bookmyshow

type PaymentGateway interface {
	Pay(amount float32, ticketID int32) (string, error)
}

type UPIPayment struct{}

func (u *UPIPayment) Pay(amount float32, ticketID int32) (string, error) {
	return "upi payment successful", nil
}

type CardPayment struct{}

func (c *CardPayment) Pay(amount float32, ticketID int32) (string, error) {
	return "card payment successful", nil
}

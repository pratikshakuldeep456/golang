package svc

type Payment interface {
	Pay(amount int32, tId int) (string, error)
}

type UPIPayment struct {
}

func (u *UPIPayment) Pay(amount int32, tId int) (string, error) {
	return "upi payment is done", nil
}

type CardPayment struct {
}

func (c *CardPayment) Pay(amount int32, tId int) (string, error) {
	return "card payment is done", nil
}

type PaymentSvc struct {
	IPayment Payment
}

func NewPaymentMethod(i Payment) *PaymentSvc {
	return &PaymentSvc{IPayment: i}
}

func (ps *PaymentSvc) MakePayment(amount int32, tId int) (string, error) {
	return ps.IPayment.Pay(amount, tId)
}

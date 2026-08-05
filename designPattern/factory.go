package designpattern

import "fmt"

type PaymentSystem interface {
	ProcessPayment(amount float64) error
}

type CreditCardPayment struct {
}

func (c *CreditCardPayment) ProcessPayment(amount float64) error {
	// Implement credit card payment processing logic here
	return nil
}

type PayPalPayment struct {
	amount float64
}

func (p *PayPalPayment) ProcessPayment(amount float64) error {
	// Implement PayPal payment processing logic here
	return nil
}

type BankTransferPayment struct {
	AccountNumber string
	IFSCCode      string
}

func (b *BankTransferPayment) ProcessPayment(amount float64) error {
	// Implement bank transfer payment processing logic here
	return nil
}

type UPI struct {
	UpiID  string
	Amount float64
}

func (u *UPI) ProcessPayment(amount float64) error {
	// Implement UPI payment processing logic here
	fmt.Println("Charging", amount, "via UPI ID:", u.UpiID)
	return nil
}

type PaymentConfig struct {
	// Add any configuration fields needed for payment processing
	Amount float64
}

func NewPaymentSystem(paymentType string, pc *PaymentConfig) PaymentSystem {
	switch paymentType {
	case "creditcard":
		return &CreditCardPayment{}
	case "paypal":
		return &PayPalPayment{}
	case "banktransfer":
		return &BankTransferPayment{}
	case "upi":
		return &UPI{}
	default:
		return nil
	}
}

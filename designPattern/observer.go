package designpattern

import "fmt"

//	type iNotification interface {
//		Notify() error
//	}
type Observer interface {
	Update(stockPrice StockPrice) error
}

type StockPrice struct {
	Symbol    string
	Price     float64
	observers []Observer
}

func (s *StockPrice) Attach(observer Observer) {
	s.observers = append(s.observers, observer)
}
func (s *StockPrice) Detach(observer Observer) {
	for _, o := range s.observers {
		o.Update(*s)
	}
}

type SMSNotifier struct {
	PhoneNumber string
}

func (s *SMSNotifier) Update(stockPrice StockPrice) error {
	// Logic to send SMS notification
	fmt.Println("Sending SMS notification to", s.PhoneNumber)
	return nil
}

type EmailNotifier struct {
	EmailAddress string
}

func (e *EmailNotifier) Update(stockPrice StockPrice) error {
	// Logic to send email notification
	//print
	fmt.Println("Sending email notification to", e.EmailAddress)
	return nil
}

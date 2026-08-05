package solid

type Accounting interface {
	GetBalance() float64
	Deposit(amount float64)
}

type SavingsAccount struct {
	Balance float64
}

func (s *SavingsAccount) GetBalance() float64 {
	return s.Balance
}

func (s *SavingsAccount) Deposit(amount float64) {
	s.Balance += amount
}

func (s *SavingsAccount) Withdraw(amount float64) {
	if amount <= s.Balance {
		s.Balance -= amount
	}
}

type FDAccount struct {
	Balance float64
}

func (f *FDAccount) GetBalance() float64 {
	return f.Balance
}

func (f *FDAccount) Deposit(amount float64) {
	f.Balance += amount
}

package fundtransfer

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// now is a variable so tests can fake the clock.
var now = time.Now

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountExists     = errors.New("account already exists")
	ErrTransferNotFound  = errors.New("transfer not found")
	ErrKeyConflict       = errors.New("idempotency key reused with different details")
	ErrInactive          = errors.New("account is frozen")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrDailyLimit        = errors.New("daily limit exceeded")
	ErrNotReversible     = errors.New("transfer cannot be reversed")
)

// ---------------- Entities ----------------
// All money is int64 paise. Never floats.

type AccountStatus int

const (
	Active AccountStatus = iota
	Frozen
)

type TransferStatus string

const (
	Initiated TransferStatus = "INITIATED"
	Completed TransferStatus = "COMPLETED"
	Failed    TransferStatus = "FAILED"
	Reversed  TransferStatus = "REVERSED"
)

type EntryType string

const (
	Debit  EntryType = "DEBIT"
	Credit EntryType = "CREDIT"
)

type Account struct {
	ID         string
	Holder     string
	Status     AccountStatus
	DailyLimit int64 // 0 = no limit

	mu       sync.Mutex // guards everything below
	balance  int64
	spent    int64
	spentDay string
	ledger   []LedgerEntry // append-only
}

func (a *Account) spentToday() int64 {
	if a.spentDay != now().Format("2006-01-02") {
		return 0
	}
	return a.spent
}

func (a *Account) addSpent(n int64) {
	day := now().Format("2006-01-02")
	if a.spentDay != day {
		a.spentDay, a.spent = day, 0
	}
	a.spent += n
}

type Transfer struct {
	ID        string
	Key       string // idempotency key
	From, To  string
	Amount    int64
	Status    TransferStatus
	Reason    string
	CreatedAt time.Time

	done      chan struct{} // closed when the transfer finishes
	err       error
	reversing bool // claim flag so two reversals cannot both run
}

type LedgerEntry struct {
	ID           int64
	TransferID   string
	AccountID    string
	Type         EntryType
	Amount       int64
	BalanceAfter int64
	Note         string // TRANSFER or REVERSAL
	At           time.Time
}

// ---------------- Rules (Chain of Responsibility) ----------------
// Called while both account locks are held, so the checks are race-free.

type Rule interface {
	Check(from, to *Account, amount int64) error
}

type ActiveRule struct{}

func (ActiveRule) Check(from, to *Account, _ int64) error {
	if from.Status != Active || to.Status != Active {
		return ErrInactive
	}
	return nil
}

type BalanceRule struct{}

func (BalanceRule) Check(from, _ *Account, amount int64) error {
	if from.balance < amount {
		return ErrInsufficientFunds
	}
	return nil
}

type DailyLimitRule struct{}

func (DailyLimitRule) Check(from, _ *Account, amount int64) error {
	if from.DailyLimit > 0 && from.spentToday()+amount > from.DailyLimit {
		return ErrDailyLimit
	}
	return nil
}

// ---------------- Service ----------------

type Service struct {
	mu       sync.Mutex // guards the maps and transfer status fields
	accounts map[string]*Account
	byID     map[string]*Transfer
	byKey    map[string]*Transfer // idempotency store
	rules    []Rule
	seq      atomic.Int64 // ids for transfers and ledger entries
}

func NewService() *Service {
	return &Service{
		accounts: map[string]*Account{},
		byID:     map[string]*Transfer{},
		byKey:    map[string]*Transfer{},
		rules:    []Rule{ActiveRule{}, BalanceRule{}, DailyLimitRule{}},
	}
}

// AddRule extends validation without touching the core. Call it before serving traffic.
func (s *Service) AddRule(r Rule) { s.rules = append(s.rules, r) }

// ---- Accounts ----

// CreateAccount: the opening balance is the only place money enters the system.
func (s *Service) CreateAccount(id, holder string, opening, dailyLimit int64) error {
	if id == "" || opening < 0 || dailyLimit < 0 {
		return ErrInvalidInput
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accounts[id]; ok {
		return ErrAccountExists
	}
	s.accounts[id] = &Account{ID: id, Holder: holder, DailyLimit: dailyLimit, balance: opening}
	return nil
}

func (s *Service) account(id string) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.accounts[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return a, nil
}

func (s *Service) GetBalance(id string) (int64, error) {
	a, err := s.account(id)
	if err != nil {
		return 0, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance, nil
}

func (s *Service) Freeze(id string) error {
	a, err := s.account(id)
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Status = Frozen
	return nil
}

func (s *Service) Statement(id string) ([]LedgerEntry, error) {
	a, err := s.account(id)
	if err != nil {
		return nil, err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]LedgerEntry(nil), a.ledger...), nil // copy: history stays immutable
}

// ---- Transfer ----

// Transfer is idempotent per key: a retry returns the original result and moves no money.
func (s *Service) Transfer(key, fromID, toID string, amount int64) (Transfer, error) {
	if key == "" || amount <= 0 || fromID == toID {
		return Transfer{}, ErrInvalidInput
	}

	s.mu.Lock()
	from, to := s.accounts[fromID], s.accounts[toID]
	if from == nil || to == nil {
		s.mu.Unlock()
		return Transfer{}, ErrAccountNotFound
	}
	// Check-then-insert on the key is atomic because s.mu is held.
	if t, ok := s.byKey[key]; ok {
		s.mu.Unlock()
		if t.From != fromID || t.To != toID || t.Amount != amount {
			return Transfer{}, ErrKeyConflict
		}
		<-t.done // a concurrent duplicate waits for the first request to finish
		return s.snapshot(t)
	}
	t := &Transfer{
		ID: fmt.Sprintf("T%d", s.seq.Add(1)), Key: key,
		From: fromID, To: toID, Amount: amount,
		Status: Initiated, CreatedAt: now(), done: make(chan struct{}),
	}
	s.byKey[key], s.byID[t.ID] = t, t
	s.mu.Unlock()

	err := s.execute(t, from, to)

	s.mu.Lock()
	if err == nil {
		t.Status = Completed
	} else {
		t.Status, t.Reason, t.err = Failed, err.Error(), err
	}
	s.mu.Unlock()
	close(t.done)
	return s.snapshot(t)
}

// execute holds both account locks for the whole check-and-move, so the debit and
// credit happen together and two transfers can never overspend the same balance.
func (s *Service) execute(t *Transfer, from, to *Account) error {
	unlock := lockPair(from, to)
	defer unlock()
	for _, r := range s.rules {
		if err := r.Check(from, to, t.Amount); err != nil {
			return err
		}
	}
	s.move(t.ID, "TRANSFER", from, to, t.Amount)
	from.addSpent(t.Amount)
	return nil
}

// move changes both balances and writes both ledger entries. Caller holds both locks.
func (s *Service) move(transferID, note string, from, to *Account, amount int64) {
	from.balance -= amount
	to.balance += amount
	at := now()
	from.ledger = append(from.ledger, LedgerEntry{s.seq.Add(1), transferID, from.ID, Debit, amount, from.balance, note, at})
	to.ledger = append(to.ledger, LedgerEntry{s.seq.Add(1), transferID, to.ID, Credit, amount, to.balance, note, at})
}

// ---- Reverse ----

// Reverse moves the money back once. Reversal does not lower the daily-limit usage.
func (s *Service) Reverse(transferID string) error {
	s.mu.Lock()
	t, ok := s.byID[transferID]
	if !ok {
		s.mu.Unlock()
		return ErrTransferNotFound
	}
	if t.Status != Completed || t.reversing { // only COMPLETED, and only one reversal at a time
		s.mu.Unlock()
		return ErrNotReversible
	}
	t.reversing = true
	from, to := s.accounts[t.From], s.accounts[t.To]
	s.mu.Unlock()

	err := s.reverseMoney(t, from, to)

	s.mu.Lock()
	t.reversing = false
	if err == nil {
		t.Status = Reversed
	}
	s.mu.Unlock()
	return err
}

func (s *Service) reverseMoney(t *Transfer, from, to *Account) error {
	unlock := lockPair(from, to)
	defer unlock()
	if to.balance < t.Amount { // receiver must still have the money
		return ErrInsufficientFunds
	}
	s.move(t.ID, "REVERSAL", to, from, t.Amount)
	return nil
}

func (s *Service) GetTransfer(id string) (Transfer, error) {
	s.mu.Lock()
	t, ok := s.byID[id]
	s.mu.Unlock()
	if !ok {
		return Transfer{}, ErrTransferNotFound
	}
	tr, _ := s.snapshot(t)
	return tr, nil
}

// ---- helpers ----

func (s *Service) snapshot(t *Transfer) (Transfer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return *t, t.err
}

// lockPair always locks the lower account ID first, so A->B and B->A can never deadlock.
func lockPair(a, b *Account) func() {
	first, second := a, b
	if second.ID < first.ID {
		first, second = second, first
	}
	first.mu.Lock()
	second.mu.Lock()
	return func() {
		second.mu.Unlock()
		first.mu.Unlock()
	}
}

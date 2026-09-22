package batchprocessing

type BatchType string

const (
	OrderBatch   BatchType = "ORDER"
	PaymentBatch BatchType = "PAYMENT"
	SignupBatch  BatchType = "SIGNUP"
)

type ItemStatus string

const (
	Pending ItemStatus = "PENDING"
	Success ItemStatus = "SUCCESS"
	Failed  ItemStatus = "FAILED"
)

type Item struct {
	ID     string
	Data   any
	Status ItemStatus
}

type Batch struct {
	ID    string
	Type  BatchType
	Items []*Item
}

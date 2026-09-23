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

//	type Item struct {
//		ID     string
//		Type   BatchType   // 👈 NEW — type now lives per-item, not just per-batch
//		Data   map[string]interface{}
//		Status Status
//	}
type Batch struct {
	ID    string
	Type  BatchType
	Items []*Item
}

package batchprocessing

// ---------- Entities ----------

// type BatchType string

// const (
// 	Order   BatchType = "ORDER"
// 	Payment BatchType = "PAYMENT"
// 	Signup  BatchType = "SIGNUP"
// )

// type Item struct {
// 	ID     string
// 	Data   any
// 	Status string
// }

// type Batch struct {
// 	ID    string
// 	Type  BatchType
// 	Items []*Item
// }

// // ---------- Chain of Responsibility ----------

// type Step interface {
// 	Execute(*Item) error
// }

// type Chain struct {
// 	step Step
// 	next *Chain
// }

// func (c *Chain) Execute(item *Item) error {
// 	if err := c.step.Execute(item); err != nil {
// 		return err
// 	}
// 	if c.next != nil {
// 		return c.next.Execute(item)
// 	}
// 	return nil
// }

// type Validate struct{}

// func (*Validate) Execute(item *Item) error {
// 	if item.Data == nil {
// 		return fmt.Errorf("invalid item")
// 	}
// 	return nil
// }

// type Transform struct{}

// func (*Transform) Execute(item *Item) error {
// 	// transform item
// 	return nil
// }

// type Save struct{}

// func (*Save) Execute(item *Item) error {
// 	// save item
// 	return nil
// }

// // ---------- Strategy ----------

// type Strategy interface {
// 	Process(*Item) error
// }

// type OrderStrategy struct{}

// func (*OrderStrategy) Process(*Item) error {
// 	// order logic
// 	return nil
// }

// type PaymentStrategy struct{}

// func (*PaymentStrategy) Process(*Item) error {
// 	// payment logic
// 	return nil
// }

// type SignupStrategy struct{}

// func (*SignupStrategy) Process(*Item) error {
// 	// signup logic
// 	return nil
// }

// // ---------- Factory ----------

// func GetStrategy(t BatchType) Strategy {
// 	switch t {
// 	case Order:
// 		return &OrderStrategy{}
// 	case Payment:
// 		return &PaymentStrategy{}
// 	case Signup:
// 		return &SignupStrategy{}
// 	}
// 	return nil
// }

// // ---------- Processor ----------

// type BatchProcessor struct {
// 	chain *Chain
// }

// func (p *BatchProcessor) Process(batch *Batch) error {

// 	strategy := GetStrategy(batch.Type)
// 	if strategy == nil {
// 		return fmt.Errorf("unsupported batch type")
// 	}

// 	for _, item := range batch.Items {

// 		// Common processing pipeline
// 		if err := p.chain.Execute(item); err != nil {
// 			item.Status = "FAILED"
// 			continue
// 		}

// 		// Batch-specific processing
// 		if err := strategy.Process(item); err != nil {
// 			item.Status = "FAILED"
// 			continue
// 		}

// 		item.Status = "SUCCESS"
// 	}

// 	return nil
// }

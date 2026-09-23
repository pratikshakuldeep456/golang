package batchprocessing

// import (
// 	"fmt"
// 	"sync"
// )

// // ============================================================
// // Models
// // ============================================================

// type ItemType string
// type Status string

// const (
// 	OrderItem   ItemType = "ORDER"
// 	PaymentItem ItemType = "PAYMENT"
// 	SignupItem  ItemType = "SIGNUP"

// 	Success Status = "SUCCESS"
// 	Failed  Status = "FAILED"
// )

// type Item struct {
// 	ID     string
// 	Type   ItemType
// 	Data   interface{}
// 	Status Status
// }

// type Batch struct {
// 	ID    string
// 	Items []*Item
// }

// // ============================================================
// // Chain of Responsibility
// // ============================================================

// type ProcessingStep interface {
// 	Execute(item *Item) error
// }

// // --------------------
// // Validation
// // --------------------

// type ValidationStep struct{}

// func (v *ValidationStep) Execute(item *Item) error {
// 	fmt.Println("Validating:", item.ID)

// 	if item.Data == nil {
// 		return fmt.Errorf("invalid item: %s", item.ID)
// 	}

// 	return nil
// }

// // --------------------
// // Transformation
// // --------------------

// type TransformationStep struct{}

// func (t *TransformationStep) Execute(item *Item) error {
// 	fmt.Println("Transforming:", item.ID)

// 	// Transformation logic

// 	return nil
// }

// // --------------------
// // Persistence
// // --------------------

// type PersistenceStep struct{}

// func (p *PersistenceStep) Execute(item *Item) error {
// 	fmt.Println("Saving:", item.ID)

// 	// DB persistence logic

// 	return nil
// }

// // --------------------
// // Chain
// // --------------------

// type ChainStep struct {
// 	step ProcessingStep
// 	next *ChainStep
// }

// func (c *ChainStep) ExecuteChain(item *Item) error {

// 	if err := c.step.Execute(item); err != nil {
// 		return err
// 	}

// 	if c.next != nil {
// 		return c.next.ExecuteChain(item)
// 	}

// 	return nil
// }

// // ============================================================
// // Strategy Pattern
// // ============================================================

// type ProcessingStrategy interface {
// 	Process(item *Item) error
// }

// // --------------------
// // Order
// // --------------------

// type OrderStrategy struct{}

// func (s *OrderStrategy) Process(item *Item) error {
// 	fmt.Println("Processing ORDER:", item.ID)

// 	// Order-specific business logic

// 	return nil
// }

// // --------------------
// // Payment
// // --------------------

// type PaymentStrategy struct{}

// func (s *PaymentStrategy) Process(item *Item) error {
// 	fmt.Println("Processing PAYMENT:", item.ID)

// 	// Payment-specific business logic

// 	return nil
// }

// // --------------------
// // Signup
// // --------------------

// type SignupStrategy struct{}

// func (s *SignupStrategy) Process(item *Item) error {
// 	fmt.Println("Processing SIGNUP:", item.ID)

// 	// Signup-specific business logic

// 	return nil
// }

// // ============================================================
// // Factory
// // ============================================================

// type StrategyFactory struct {
// 	strategies map[ItemType]ProcessingStrategy
// }

// func NewStrategyFactory() *StrategyFactory {

// 	return &StrategyFactory{
// 		strategies: map[ItemType]ProcessingStrategy{
// 			OrderItem:   &OrderStrategy{},
// 			PaymentItem: &PaymentStrategy{},
// 			SignupItem:  &SignupStrategy{},
// 		},
// 	}
// }

// func (f *StrategyFactory) GetStrategy(
// 	itemType ItemType,
// ) ProcessingStrategy {

// 	return f.strategies[itemType]
// }

// // ============================================================
// // Batch Processor
// // ============================================================

// type BatchProcessor struct {
// 	factory    *StrategyFactory
// 	chain      *ChainStep
// 	workerSize int
// }

// // ============================================================
// // Concurrent Batch Processing
// // ============================================================

// func (p *BatchProcessor) BatchProcess(batch *Batch) error {

// 	jobs := make(chan *Item)

// 	var wg sync.WaitGroup

// 	// --------------------------------------------------------
// 	// Create workers
// 	// --------------------------------------------------------

// 	for i := 0; i < p.workerSize; i++ {

// 		wg.Add(1)

// 		go func(workerID int) {

// 			defer wg.Done()

// 			for item := range jobs {

// 				fmt.Printf(
// 					"Worker %d processing %s\n",
// 					workerID,
// 					item.ID,
// 				)

// 				// ------------------------------------------------
// 				// 1. Common processing
// 				// ------------------------------------------------

// 				if err := p.chain.ExecuteChain(item); err != nil {

// 					item.Status = Failed

// 					fmt.Printf(
// 						"Item %s failed in common processing: %v\n",
// 						item.ID,
// 						err,
// 					)

// 					continue
// 				}

// 				// ------------------------------------------------
// 				// 2. Get strategy based on ITEM TYPE
// 				// ------------------------------------------------

// 				strategy := p.factory.GetStrategy(item.Type)

// 				if strategy == nil {

// 					item.Status = Failed

// 					fmt.Printf(
// 						"Unsupported item type: %s\n",
// 						item.Type,
// 					)

// 					continue
// 				}

// 				// ------------------------------------------------
// 				// 3. Type-specific processing
// 				// ------------------------------------------------

// 				if err := strategy.Process(item); err != nil {

// 					item.Status = Failed

// 					fmt.Printf(
// 						"Item %s failed in strategy: %v\n",
// 						item.ID,
// 						err,
// 					)

// 					continue
// 				}

// 				// ------------------------------------------------
// 				// 4. Success
// 				// ------------------------------------------------

// 				item.Status = Success

// 				fmt.Printf(
// 					"Item %s processed successfully\n",
// 					item.ID,
// 				)
// 			}

// 		}(i)
// 	}

// 	// --------------------------------------------------------
// 	// Send items to workers
// 	// --------------------------------------------------------

// 	for _, item := range batch.Items {
// 		jobs <- item
// 	}

// 	close(jobs)

// 	// --------------------------------------------------------
// 	// Wait for workers
// 	// --------------------------------------------------------

// 	wg.Wait()

// 	return nil
// }

// // ============================================================
// // Constructor
// // ============================================================

// func NewBatchProcessor(workerSize int) *BatchProcessor {

// 	// Build validation step
// 	validation := &ChainStep{
// 		step: &ValidationStep{},
// 	}

// 	// Build transformation step
// 	transformation := &ChainStep{
// 		step: &TransformationStep{},
// 	}

// 	// Build persistence step
// 	persistence := &ChainStep{
// 		step: &PersistenceStep{},
// 	}

// 	// Validation
// 	//      ↓
// 	// Transformation
// 	//      ↓
// 	// Persistence

// 	validation.next = transformation
// 	transformation.next = persistence

// 	return &BatchProcessor{
// 		factory:    NewStrategyFactory(),
// 		chain:      validation,
// 		workerSize: workerSize,
// 	}
// }

// // ============================================================
// // Main
// // ============================================================

// func main() {

// 	batch := &Batch{
// 		ID: "BATCH-001",
// 		Items: []*Item{

// 			{
// 				ID:   "ORDER-1",
// 				Type: OrderItem,
// 				Data: "Order Data",
// 			},

// 			{
// 				ID:   "PAYMENT-1",
// 				Type: PaymentItem,
// 				Data: "Payment Data",
// 			},

// 			{
// 				ID:   "SIGNUP-1",
// 				Type: SignupItem,
// 				Data: "Signup Data",
// 			},

// 			{
// 				ID:   "ORDER-2",
// 				Type: OrderItem,
// 				Data: "Order Data",
// 			},

// 			{
// 				ID:   "PAYMENT-2",
// 				Type: PaymentItem,
// 				Data: "Payment Data",
// 			},

// 			{
// 				ID:   "INVALID-1",
// 				Type: OrderItem,
// 				Data: nil,
// 			},
// 		},
// 	}

// 	// Create processor with 3 workers
// 	processor := NewBatchProcessor(3)

// 	// Process batch
// 	if err := processor.BatchProcess(batch); err != nil {
// 		fmt.Println("Batch failed:", err)
// 	}

// 	// --------------------------------------------------------
// 	// Print final status
// 	// --------------------------------------------------------

// 	fmt.Println("\nFinal Status:")

// 	for _, item := range batch.Items {

// 		fmt.Printf(
// 			"%s (%s) -> %s\n",
// 			item.ID,
// 			item.Type,
// 			item.Status,
// 		)
// 	}
// }

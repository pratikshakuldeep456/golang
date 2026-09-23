package batchprocessing

import "fmt"

// ====================
// Chain of Responsibility
// ====================

type ProcessingStep interface {
	Execute(item *Item) error
}

type ValidationStep struct{}

func (v *ValidationStep) Execute(item *Item) error {
	fmt.Println("Validating:", item.ID)

	if item.Data == nil {
		return fmt.Errorf("invalid item: %s", item.ID)
	}

	return nil
}

type TransformationStep struct{}

func (t *TransformationStep) Execute(item *Item) error {
	fmt.Println("Transforming:", item.ID)

	// transformation logic

	return nil
}

type PersistenceStep struct{}

func (p *PersistenceStep) Execute(item *Item) error {
	fmt.Println("Saving:", item.ID)

	// save to DB

	return nil
}

type ChainStep struct {
	step ProcessingStep
	next *ChainStep
}

func (c *ChainStep) ExecuteChain(item *Item) error {

	if err := c.step.Execute(item); err != nil {
		return err
	}

	if c.next != nil {
		return c.next.ExecuteChain(item)
	}

	return nil
}

// ====================
// Strategy
// ====================

type ProcessingStrategy interface {
	Process(item *Item) error
}

type OrderStrategy struct{}

func (s *OrderStrategy) Process(item *Item) error {
	fmt.Println("Processing ORDER:", item.ID)

	// order-specific business logic

	return nil
}

type PaymentStrategy struct{}

func (s *PaymentStrategy) Process(item *Item) error {
	fmt.Println("Processing PAYMENT:", item.ID)

	// payment-specific business logic

	return nil
}

type SignupStrategy struct{}

func (s *SignupStrategy) Process(item *Item) error {
	fmt.Println("Processing SIGNUP:", item.ID)

	// signup-specific business logic

	return nil
}

// ====================
// Factory
// ====================

type StrategyFactory struct{}

func (f *StrategyFactory) GetStrategy(
	batchType BatchType,
) ProcessingStrategy {

	switch batchType {

	case OrderBatch:
		return &OrderStrategy{}

	case PaymentBatch:
		return &PaymentStrategy{}

	case SignupBatch:
		return &SignupStrategy{}

	default:
		return nil
	}
}

// ====================
// Batch Processor
// ====================
/*BatchProcessor is the orchestrator —
the "brain" that owns and coordinates
the two patterns you just went through
(Factory + Chain of Responsibility).
It doesn't do any actual
validation/transformation/business-logic itself
 — it just holds references to the things that do,
 and calls them in the right order.*/

type BatchProcessor struct {
	factory *StrategyFactory
	chain   *ChainStep
	//workerSize int
}

//concurrent
/*
func (p *BatchProcessor) Process(batch *Batch) error {

	strategy := p.factory.GetStrategy(batch.Type)

	if strategy == nil {
		return fmt.Errorf("unsupported batch type: %s", batch.Type)
	}

	// Jobs channel
	jobs := make(chan *Item)

	var wg sync.WaitGroup

	// Create workers
	for i := 0; i < p.workerSize; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()

			for item := range jobs {

				// Steps are sequential for one item
				if err := p.chain.Execute(item); err != nil {
					item.Status = Failed
					fmt.Println("Item failed:", item.ID)
					continue
				}

				// Batch-type-specific processing
				if err := strategy.Process(item); err != nil {
					item.Status = Failed
					fmt.Println("Item failed:", item.ID)
					continue
				}

				item.Status = Success
			}
		}()
	}

	// Send items to workers
	for _, item := range batch.Items {
		jobs <- item
	}

	close(jobs)

	// Wait for all workers
	wg.Wait()

	return nil
} */

func (p *BatchProcessor) BatchProcess(batch *Batch) error {

	//re out which business logic applies to
	// this whole batch (e.g., is this an Order batch or a Payment batch?).
	strategy := p.factory.GetStrategy(batch.Type)
	if strategy == nil {
		return fmt.Errorf("unsupported batch type: %s", batch.Type)
	}
	for _, item := range batch.Items {
		//   strategy := p.factory.GetStrategy(batch.Type)
		// 	if strategy == nil {
		// 		return fmt.Errorf("unsupported batch type: %s", batch.Type)
		// 	}
		// First execute common processing steps
		//Run this one item through the Chain of Responsibility
		// (Validation → Transformation → Persistence).
		if err := p.chain.ExecuteChain(item); err != nil {
			item.Status = Failed
			fmt.Println("Item failed:", item.ID)
			continue
		}

		// Then execute batch-type-specific logic
		if err := strategy.Process(item); err != nil {
			item.Status = Failed
			fmt.Println("Item failed:", item.ID)
			continue
		}

		item.Status = Success
	}

	return nil
}

// ====================
// Build Processor
// ====================

// ValidationStep     → validates item
// TransformationStep → transforms item
// PersistenceStep    → saves item
func NewBatchProcessor() *BatchProcessor {

	validation := &ChainStep{
		step: &ValidationStep{},
	}

	transformation := &ChainStep{
		step: &TransformationStep{},
	}

	persistence := &ChainStep{
		step: &PersistenceStep{},
	}

	// Build chain:
	// Validation → Transformation → Persistence

	validation.next = transformation
	transformation.next = persistence

	return &BatchProcessor{
		factory: &StrategyFactory{},
		chain:   validation,
	}
}

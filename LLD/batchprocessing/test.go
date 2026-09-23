package batchprocessing

func TestBatchProcessor() {
	processorInstance := NewBatchProcessor()

	Batchs := &Batch{
		ID:   "12",
		Type: OrderBatch,
		Items: []*Item{
			{ID: "order-1", Data: map[string]interface{}{"amount": 500}},
			{ID: "order-2", Data: map[string]interface{}{"amount": 1200}},
		},
	}

	processorInstance.BatchProcess(Batchs)
}

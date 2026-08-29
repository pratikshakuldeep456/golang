package asyncworker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ---------- Data model ----------

type Task struct {
	ID      string
	Ctx     context.Context
	Execute func(ctx context.Context) (any, error)
}

type Result struct {
	TaskID string
	Value  any
	Err    error
}

// ---------- Worker Pool ----------

type WorkerPool struct {
	tasks   chan Task
	results chan Result
	wg      sync.WaitGroup
	quit    chan struct{}
	once    sync.Once
}

func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
	p := &WorkerPool{
		tasks:   make(chan Task, queueSize),
		results: make(chan Result, queueSize),
		quit:    make(chan struct{}),
	}
	for i := 0; i < numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	return p
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			p.execute(task)
		case <-p.quit:
			return
		}
	}
}

func (p *WorkerPool) execute(task Task) {
	ctx := task.Ctx
	done := make(chan Result, 1)

	go func() {
		val, err := task.Execute(ctx)
		done <- Result{TaskID: task.ID, Value: val, Err: err}
	}()

	select {
	case res := <-done:
		p.results <- res
	case <-ctx.Done():
		p.results <- Result{TaskID: task.ID, Err: ctx.Err()}
	}
}

func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.tasks <- task:
		return nil
	case <-task.Ctx.Done():
		return task.Ctx.Err()
	default:
		return errors.New("queue full")
	}
}

func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

func (p *WorkerPool) Shutdown(ctx context.Context) error {
	p.once.Do(func() { close(p.tasks) })

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		close(p.quit)
		return ctx.Err()
	}
}

// ---------- Demo ----------

func AsyncWorker() {
	pool := NewWorkerPool(3, 100) // 3 workers, queue of 10

	// A slow task that respects context (finishes fine)
	makeTask := func(id string, work time.Duration) Task {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = cancel // in real code, store & call cancel() after task completes
		return Task{
			ID:  id,
			Ctx: ctx,
			Execute: func(ctx context.Context) (any, error) {
				select {
				case <-time.After(work):
					return fmt.Sprintf("task %s done", id), nil
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			},
		}
	}

	// Submit a few tasks: some finish fine, one times out
	tasksToSubmit := []Task{
		makeTask("t1", 500*time.Millisecond),
		makeTask("t2", 300*time.Millisecond),
		makeTask("t3", 3*time.Second), // exceeds its 2s timeout -> will be canceled
		makeTask("t4", 1*time.Second), // exceeds its 2s timeout -> will be canceled
	}

	for _, t := range tasksToSubmit {
		if err := pool.Submit(t); err != nil {
			fmt.Println("submit error:", err)
		}
	}

	// Collect results as they come in
	go func() {
		for i := 0; i < len(tasksToSubmit); i++ {
			res := <-pool.Results()
			if res.Err != nil {
				fmt.Printf("[%s] failed: %v\n", res.TaskID, res.Err)
			} else {
				fmt.Printf("[%s] result: %v\n", res.TaskID, res.Value)
			}
		}
	}()

	// Give tasks time to run, then shut down gracefully
	time.Sleep(4 * time.Second)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Shutdown(shutdownCtx); err != nil {
		fmt.Println("shutdown error:", err)
	} else {
		fmt.Println("pool shut down cleanly")
	}
}

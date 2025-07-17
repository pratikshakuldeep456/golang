package notes

import (
	"context"
	"fmt"
	"time"
)

func CtxC(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func CtxB(ctx context.Context) error {
	return CtxC(ctx)
}

func CtxA(ctx context.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()

	err := CtxB(ctx)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func ContextFunc() {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	cancel()

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Finished task")
	case <-ctx.Done():
		fmt.Println("Error:", ctx.Err()) // Triggers after 1 second
	}
	// CtxA(ctx)
}

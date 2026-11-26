package engine

import (
	"context"
	"sync"
	"testing"
)

func TestBatchDispatcher(t *testing.T) {
	sink := NewMemorySink()
	disp := NewBatchDispatcher(sink, 3)

	ctx := context.Background()

	// Submit 2 events - should remain buffered
	e1 := TransitionEvent{RecordID: "r1", Transition: TransitionEntered}
	e2 := TransitionEvent{RecordID: "r2", Transition: TransitionEntered}

	if err := disp.Submit(ctx, e1, e2); err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if len(sink.Events()) != 0 {
		t.Fatalf("expected 0 events in sink before batch threshold, got %d", len(sink.Events()))
	}

	// Submit 3rd event - triggers batch write
	e3 := TransitionEvent{RecordID: "r3", Transition: TransitionEntered}
	if err := disp.Submit(ctx, e3); err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if len(sink.Events()) != 3 {
		t.Fatalf("expected 3 events in sink after batch threshold, got %d", len(sink.Events()))
	}

	// Submit 4th event and manual flush
	e4 := TransitionEvent{RecordID: "r4", Transition: TransitionExited}
	if err := disp.Submit(ctx, e4); err != nil {
		t.Fatalf("submit failed: %v", err)
	}
	if err := disp.Flush(ctx); err != nil {
		t.Fatalf("flush failed: %v", err)
	}
	if len(sink.Events()) != 4 {
		t.Fatalf("expected 4 events in sink after flush, got %d", len(sink.Events()))
	}
}

func TestBatchDispatcherConcurrency(t *testing.T) {
	sink := NewMemorySink()
	disp := NewBatchDispatcher(sink, 10)
	ctx := context.Background()

	var wg sync.WaitGroup
	workers := 10
	itemsPerWorker := 50

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < itemsPerWorker; j++ {
				_ = disp.Submit(ctx, TransitionEvent{
					RecordID:   "rec",
					Transition: TransitionEntered,
				})
			}
		}(i)
	}
	wg.Wait()
	_ = disp.Flush(ctx)

	total := len(sink.Events())
	expected := workers * itemsPerWorker
	if total != expected {
		t.Fatalf("expected %d total events, got %d", expected, total)
	}
}

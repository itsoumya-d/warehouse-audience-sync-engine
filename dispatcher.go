package engine

import (
	"context"
	"sync"
)

// Sink represents a downstream sync destination (e.g. CRM, email API, ad platform).
type Sink interface {
	WriteBatch(ctx context.Context, events []TransitionEvent) error
}

// MemorySink is an in-memory sink implementation for testing and validation.
type MemorySink struct {
	mu     sync.Mutex
	events []TransitionEvent
}

func NewMemorySink() *MemorySink {
	return &MemorySink{events: make([]TransitionEvent, 0)}
}

func (m *MemorySink) WriteBatch(_ context.Context, events []TransitionEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, events...)
	return nil
}

func (m *MemorySink) Events() []TransitionEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	cpy := make([]TransitionEvent, len(m.events))
	copy(cpy, m.events)
	return cpy
}

// BatchDispatcher buffers transitions and writes to sinks when batch size is reached.
type BatchDispatcher struct {
	sink      Sink
	batchSize int
	mu        sync.Mutex
	buffer    []TransitionEvent
}

func NewBatchDispatcher(sink Sink, batchSize int) *BatchDispatcher {
	if batchSize <= 0 {
		batchSize = 100
	}
	return &BatchDispatcher{
		sink:      sink,
		batchSize: batchSize,
		buffer:    make([]TransitionEvent, 0, batchSize),
	}
}

func (d *BatchDispatcher) Submit(ctx context.Context, events ...TransitionEvent) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.buffer = append(d.buffer, events...)
	if len(d.buffer) >= d.batchSize {
		return d.flushLocked(ctx)
	}
	return nil
}

func (d *BatchDispatcher) Flush(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.flushLocked(ctx)
}

func (d *BatchDispatcher) flushLocked(ctx context.Context) error {
	if len(d.buffer) == 0 {
		return nil
	}
	toSend := make([]TransitionEvent, len(d.buffer))
	copy(toSend, d.buffer)
	d.buffer = d.buffer[:0]
	return d.sink.WriteBatch(ctx, toSend)
}

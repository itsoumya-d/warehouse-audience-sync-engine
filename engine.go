package engine

import "errors"

// OpType represents the CDC operation type.
type OpType string

const (
	OpInsert OpType = "INSERT"
	OpUpdate OpType = "UPDATE"
	OpDelete OpType = "DELETE"
)

// Record represents a customer record state from a data warehouse.
type Record struct {
	ID         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// CDCEvent represents an incoming change data capture event from a warehouse stream.
type CDCEvent struct {
	Op       OpType  `json:"op"`
	Before   *Record `json:"before,omitempty"`
	After    *Record `json:"after,omitempty"`
	Sequence uint64  `json:"sequence"`
}

// ErrInvalidEvent is returned when an event lacks required records.
var ErrInvalidEvent = errors.New("invalid cdc event")

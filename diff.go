package engine

import "time"

// TransitionType represents an audience membership state change.
type TransitionType string

const (
	TransitionEntered    TransitionType = "ENTERED"
	TransitionExited     TransitionType = "EXITED"
	TransitionMaintained TransitionType = "MAINTAINED"
	TransitionNone       TransitionType = "NONE"
)

// TransitionEvent captures an audience change resulting from a CDC operation.
type TransitionEvent struct {
	SegmentID    string         `json:"segment_id"`
	RecordID     string         `json:"record_id"`
	Transition   TransitionType `json:"transition"`
	EventOp      OpType         `json:"event_op"`
	Sequence     uint64         `json:"sequence"`
	Attributes   map[string]any `json:"attributes"`
	EvaluatedAt  time.Time      `json:"evaluated_at"`
}

// DiffEngine processes CDC events against active audience segments.
type DiffEngine struct {
	segments map[string]Segment
}

// NewDiffEngine creates a new DiffEngine.
func NewDiffEngine(segments ...Segment) *DiffEngine {
	m := make(map[string]Segment, len(segments))
	for _, s := range segments {
		m[s.ID] = s
	}
	return &DiffEngine{segments: m}
}

// ProcessCDC processes a single CDC event and returns all actionable audience transitions.
func (e *DiffEngine) ProcessCDC(event CDCEvent) ([]TransitionEvent, error) {
	if event.Op == OpInsert && event.After == nil {
		return nil, ErrInvalidEvent
	}
	if event.Op == OpDelete && event.Before == nil {
		return nil, ErrInvalidEvent
	}
	if event.Op == OpUpdate && (event.Before == nil || event.After == nil) {
		return nil, ErrInvalidEvent
	}

	recordID := ""
	if event.After != nil {
		recordID = event.After.ID
	} else if event.Before != nil {
		recordID = event.Before.ID
	}

	now := time.Now().UTC()
	var transitions []TransitionEvent

	for _, seg := range e.segments {
		beforeMatches := seg.Matches(event.Before)
		afterMatches := seg.Matches(event.After)

		var tr TransitionType
		switch {
		case !beforeMatches && afterMatches:
			tr = TransitionEntered
		case beforeMatches && !afterMatches:
			tr = TransitionExited
		case beforeMatches && afterMatches:
			tr = TransitionMaintained
		default:
			tr = TransitionNone
		}

		if tr != TransitionNone {
			attrs := make(map[string]any)
			if event.After != nil {
				for k, v := range event.After.Attributes {
					attrs[k] = v
				}
			} else if event.Before != nil {
				for k, v := range event.Before.Attributes {
					attrs[k] = v
				}
			}
			transitions = append(transitions, TransitionEvent{
				SegmentID:   seg.ID,
				RecordID:    recordID,
				Transition:  tr,
				EventOp:     event.Op,
				Sequence:    event.Sequence,
				Attributes:  attrs,
				EvaluatedAt: now,
			})
		}
	}

	return transitions, nil
}

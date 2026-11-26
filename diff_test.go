package engine

import "testing"

func TestDiffEngineInsert(t *testing.T) {
	seg := Segment{
		ID:    "active-customers",
		Logic: LogicAnd,
		Predicates: []Predicate{
			{Field: "status", Operator: OpEq, Value: "active"},
		},
	}
	engine := NewDiffEngine(seg)

	evt := CDCEvent{
		Op: OpInsert,
		After: &Record{
			ID: "cust-100",
			Attributes: map[string]interface{}{
				"status": "active",
			},
		},
		Sequence: 1,
	}

	transitions, err := engine.ProcessCDC(evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(transitions) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(transitions))
	}
	if transitions[0].Transition != TransitionEntered {
		t.Errorf("expected ENTERED, got %s", transitions[0].Transition)
	}
	if transitions[0].RecordID != "cust-100" {
		t.Errorf("expected cust-100, got %s", transitions[0].RecordID)
	}
}

func TestDiffEngineUpdateTransitions(t *testing.T) {
	seg := Segment{
		ID:    "vip",
		Logic: LogicAnd,
		Predicates: []Predicate{
			{Field: "score", Operator: OpGte, Value: 100},
		},
	}
	engine := NewDiffEngine(seg)

	// Case 1: Qualifies (Entered)
	evt1 := CDCEvent{
		Op:       OpUpdate,
		Before:   &Record{ID: "c1", Attributes: map[string]interface{}{"score": 50}},
		After:    &Record{ID: "c1", Attributes: map[string]interface{}{"score": 120}},
		Sequence: 10,
	}
	tr1, _ := engine.ProcessCDC(evt1)
	if len(tr1) != 1 || tr1[0].Transition != TransitionEntered {
		t.Fatalf("expected ENTERED, got %v", tr1)
	}

	// Case 2: Disqualifies (Exited)
	evt2 := CDCEvent{
		Op:       OpUpdate,
		Before:   &Record{ID: "c1", Attributes: map[string]interface{}{"score": 120}},
		After:    &Record{ID: "c1", Attributes: map[string]interface{}{"score": 40}},
		Sequence: 11,
	}
	tr2, _ := engine.ProcessCDC(evt2)
	if len(tr2) != 1 || tr2[0].Transition != TransitionExited {
		t.Fatalf("expected EXITED, got %v", tr2)
	}

	// Case 3: Retains qualification (Maintained)
	evt3 := CDCEvent{
		Op:       OpUpdate,
		Before:   &Record{ID: "c1", Attributes: map[string]interface{}{"score": 120}},
		After:    &Record{ID: "c1", Attributes: map[string]interface{}{"score": 150}},
		Sequence: 12,
	}
	tr3, _ := engine.ProcessCDC(evt3)
	if len(tr3) != 1 || tr3[0].Transition != TransitionMaintained {
		t.Fatalf("expected MAINTAINED, got %v", tr3)
	}

	// Case 4: Neither before nor after meets criteria (None)
	evt4 := CDCEvent{
		Op:       OpUpdate,
		Before:   &Record{ID: "c1", Attributes: map[string]interface{}{"score": 20}},
		After:    &Record{ID: "c1", Attributes: map[string]interface{}{"score": 30}},
		Sequence: 13,
	}
	tr4, _ := engine.ProcessCDC(evt4)
	if len(tr4) != 0 {
		t.Fatalf("expected 0 transitions, got %d", len(tr4))
	}
}

func TestDiffEngineDelete(t *testing.T) {
	seg := Segment{
		ID:    "active",
		Logic: LogicAnd,
		Predicates: []Predicate{
			{Field: "active", Operator: OpEq, Value: "true"},
		},
	}
	engine := NewDiffEngine(seg)

	evt := CDCEvent{
		Op:       OpDelete,
		Before:   &Record{ID: "c-del", Attributes: map[string]interface{}{"active": "true"}},
		Sequence: 99,
	}
	tr, err := engine.ProcessCDC(evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tr) != 1 || tr[0].Transition != TransitionExited {
		t.Fatalf("expected EXITED, got %v", tr)
	}
}

func TestDiffEngineInvalidEvents(t *testing.T) {
	engine := NewDiffEngine()

	if _, err := engine.ProcessCDC(CDCEvent{Op: OpInsert}); err != ErrInvalidEvent {
		t.Errorf("expected ErrInvalidEvent on empty insert")
	}
	if _, err := engine.ProcessCDC(CDCEvent{Op: OpDelete}); err != ErrInvalidEvent {
		t.Errorf("expected ErrInvalidEvent on empty delete")
	}
	if _, err := engine.ProcessCDC(CDCEvent{Op: OpUpdate, Before: &Record{}}); err != ErrInvalidEvent {
		t.Errorf("expected ErrInvalidEvent on update missing after")
	}
}

package engine

import "testing"

func TestEvaluatorAND(t *testing.T) {
	seg := Segment{
		ID:    "high-spenders",
		Name:  "High Spenders Tier",
		Logic: LogicAnd,
		Predicates: []Predicate{
			{Field: "lifetime_spend", Operator: OpGte, Value: 1000.0},
			{Field: "plan", Operator: OpEq, Value: "enterprise"},
		},
	}

	recPass := &Record{
		ID: "u-1",
		Attributes: map[string]interface{}{
			"lifetime_spend": 1500.0,
			"plan":           "enterprise",
		},
	}
	if !seg.Matches(recPass) {
		t.Fatalf("expected recPass to match high-spenders")
	}

	recFailSpend := &Record{
		ID: "u-2",
		Attributes: map[string]interface{}{
			"lifetime_spend": 800.0,
			"plan":           "enterprise",
		},
	}
	if seg.Matches(recFailSpend) {
		t.Fatalf("expected recFailSpend to NOT match high-spenders")
	}

	recFailPlan := &Record{
		ID: "u-3",
		Attributes: map[string]interface{}{
			"lifetime_spend": 2500.0,
			"plan":           "starter",
		},
	}
	if seg.Matches(recFailPlan) {
		t.Fatalf("expected recFailPlan to NOT match high-spenders")
	}
}

func TestEvaluatorOR(t *testing.T) {
	seg := Segment{
		ID:    "engaged-or-enterprise",
		Name:  "Engaged or Enterprise",
		Logic: LogicOr,
		Predicates: []Predicate{
			{Field: "login_count", Operator: OpGt, Value: 50},
			{Field: "domain", Operator: OpContains, Value: "acme.corp"},
		},
	}

	r1 := &Record{ID: "u-1", Attributes: map[string]interface{}{"login_count": 60, "domain": "other.com"}}
	r2 := &Record{ID: "u-2", Attributes: map[string]interface{}{"login_count": 10, "domain": "lead@acme.corp"}}
	r3 := &Record{ID: "u-3", Attributes: map[string]interface{}{"login_count": 5, "domain": "unknown.com"}}

	if !seg.Matches(r1) {
		t.Errorf("r1 should match")
	}
	if !seg.Matches(r2) {
		t.Errorf("r2 should match")
	}
	if seg.Matches(r3) {
		t.Errorf("r3 should not match")
	}
}

func TestEvaluatorEdgeCases(t *testing.T) {
	seg := Segment{
		ID:    "test",
		Logic: LogicAnd,
		Predicates: []Predicate{
			{Field: "missing_field", Operator: OpEq, Value: "test"},
		},
	}
	r := &Record{ID: "u-1", Attributes: map[string]interface{}{"other": 123}}
	if seg.Matches(r) {
		t.Errorf("missing field should return false")
	}
	if seg.Matches(nil) {
		t.Errorf("nil record should return false")
	}
}

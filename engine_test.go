package engine

import "testing"

func TestBaselineRecord(t *testing.T) {
	rec := &Record{
		ID: "cust-1",
		Attributes: map[string]interface{}{
			"spend": 120.5,
		},
	}
	if rec.ID != "cust-1" {
		t.Fatalf("expected ID cust-1, got %s", rec.ID)
	}
}

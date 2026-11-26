package engine

import "testing"

func BenchmarkProcessCDC(b *testing.B) {
	seg := Segment{
		ID:    "benchmark-segment",
		Logic: LogicAnd,
		Predicates: []Predicate{
			{Field: "mrr", Operator: OpGt, Value: 500.0},
			{Field: "active", Operator: OpEq, Value: "true"},
		},
	}
	diffEngine := NewDiffEngine(seg)

	evt := CDCEvent{
		Op:       OpUpdate,
		Before:   &Record{ID: "bench-user", Attributes: map[string]interface{}{"mrr": 400.0, "active": "true"}},
		After:    &Record{ID: "bench-user", Attributes: map[string]interface{}{"mrr": 650.0, "active": "true"}},
		Sequence: 1001,
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = diffEngine.ProcessCDC(evt)
	}
}

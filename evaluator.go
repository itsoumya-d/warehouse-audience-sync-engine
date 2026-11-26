package engine

import (
	"fmt"
	"strings"
)

type PredicateOp string

const (
	OpEq       PredicateOp = "eq"
	OpNeq      PredicateOp = "neq"
	OpGt       PredicateOp = "gt"
	OpGte      PredicateOp = "gte"
	OpLt       PredicateOp = "lt"
	OpLte      PredicateOp = "lte"
	OpContains PredicateOp = "contains"
)

// Predicate defines a condition against a record attribute.
type Predicate struct {
	Field    string      `json:"field"`
	Operator PredicateOp `json:"operator"`
	Value    interface{} `json:"value"`
}

// LogicType combines multiple predicates.
type LogicType string

const (
	LogicAnd LogicType = "AND"
	LogicOr  LogicType = "OR"
)

// Segment defines a target audience cohort.
type Segment struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Logic      LogicType   `json:"logic"`
	Predicates []Predicate `json:"predicates"`
}

// Matches evaluates whether a record qualifies for a segment.
func (s Segment) Matches(rec *Record) bool {
	if rec == nil || len(s.Predicates) == 0 {
		return false
	}
	if s.Logic == LogicOr {
		for _, p := range s.Predicates {
			if evalPredicate(rec, p) {
				return true
			}
		}
		return false
	}
	// Default to AND logic
	for _, p := range s.Predicates {
		if !evalPredicate(rec, p) {
			return false
		}
	}
	return true
}

func evalPredicate(rec *Record, p Predicate) bool {
	val, ok := rec.Attributes[p.Field]
	if !ok || val == nil {
		return false
	}

	switch p.Operator {
	case OpEq:
		return fmt.Sprintf("%v", val) == fmt.Sprintf("%v", p.Value)
	case OpNeq:
		return fmt.Sprintf("%v", val) != fmt.Sprintf("%v", p.Value)
	case OpContains:
		return strings.Contains(strings.ToLower(fmt.Sprintf("%v", val)), strings.ToLower(fmt.Sprintf("%v", p.Value)))
	case OpGt, OpGte, OpLt, OpLte:
		numVal, okVal := toFloat(val)
		target, okTarget := toFloat(p.Value)
		if !okVal || !okTarget {
			return false
		}
		switch p.Operator {
		case OpGt:
			return numVal > target
		case OpGte:
			return numVal >= target
		case OpLt:
			return numVal < target
		case OpLte:
			return numVal <= target
		}
	}
	return false
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	default:
		return 0, false
	}
}

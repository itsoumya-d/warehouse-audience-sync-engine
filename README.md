# Warehouse Audience Sync Engine

A high-throughput CDC (Change Data Capture) stream consumer and real-time audience segment evaluation engine for warehouse data activation and reverse ETL workflows.

## Features
- **CDC Ingestion & State Diffing**: Ingests change streams (Insert, Update, Delete) from analytical data warehouses, evaluating before/after state transitions (`ENTERED`, `EXITED`, `MAINTAINED`).
- **Dynamic Segment Evaluation**: Evaluates rule-based predicates (numerical thresholds, categorical matches, temporal activity bounds) against customer traits with sub-microsecond latency.
- **Downstream Sync Dispatcher**: Buffers and dispatches audience transition events to external activation targets with batched delivery.

## Running Tests & Benchmarks
```bash
go test -v -race -cover ./...
go test -bench=. -benchmem ./...
```

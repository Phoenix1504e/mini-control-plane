# MVCC Conflict Scaling Campaign

## Purpose

This experimental campaign evaluates how MVCC conflict frequency scales
under increasing reconciliation concurrency in a mini control-plane
implementation backed by etcd.

The study focuses on:
- optimistic concurrency contention
- stale reconciliation behavior
- compound fault interaction
- reconciliation saturation regimes

---

# Experimental Configuration

## Datastore

- etcd v3
- persistent datastore across trials
- resources reset between runs

---

## Controller Configuration

Controllers tested:

- 1 controller
- 2 controllers
- 4 controllers
- 8 controllers

Each configuration was executed across multiple repeated trials.

---

# Workload Parameters

| Parameter | Value |
|---|---|
| Resources | 10 |
| Experiment Duration | 20s |
| Reconcile Delay | 1s |
| Trials Per Configuration | 5 |

---

# Fault Injection

The following perturbations were enabled:

- watch-event loss
- delayed reconciliation
- concurrent status updates

Watch perturbation simulates informer inconsistency and stale state
propagation.

---

# Metrics Collected

Telemetry captured:

- MVCC conflict events
- reconcile latency
- reconcile failures
- state samples
- leader election events

---

# Hypothesis

MVCC conflict frequency increases as reconciliation concurrency grows.

Formally:

C(N) increases with concurrent reconcilers N.

---

# Experimental Notes

Earlier exploratory campaigns experienced:
- datastore saturation
- orchestration leakage
- WSL instability
- WAL teardown races

This campaign isolates experiments using:
- persistent etcd lifecycle
- ephemeral controllers
- watchdog-based orchestration
- per-trial resource cleanup

---

# Limitations

Current limitations include:

- single-node datastore deployment
- local WSL execution environment
- synthetic workloads
- limited resource cardinality

Results should therefore be interpreted as exploratory systems behavior
rather than production-scale benchmarking.

---

# Directory Structure

Each run directory contains:
- metadata.json
- telemetry artifacts
- experimental outputs

Run naming follows:

run-<timestamp>

---

# Reproducibility

Experimental campaigns can be reproduced using:

scripts/run_mvcc_trials.sh

Aggregation:

analysis/aggregate_conflicts.py

Plot generation:

analysis/plot_conflicts.py

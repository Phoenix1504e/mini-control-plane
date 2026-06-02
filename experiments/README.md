# Experiments

This directory documents the fault-injection experiments used to evaluate Mini Control Plane's white-box control-plane semantics.

The Table 1 configurations are versioned under `specs/experiments/`. Each run varies the number of competing controllers while keeping the lease TTL and retry interval fixed:

| Config | Controllers | Lease TTL | Retry Interval |
|--------|-------------|-----------|----------------|
| `table1-n1.yaml` | 1 | 5s | 2s |
| `table1-n2.yaml` | 2 | 5s | 2s |
| `table1-n4.yaml` | 4 | 5s | 2s |
| `table1-n8.yaml` | 8 | 5s | 2s |

Structured JSONL logs should be written to `experiments/results/` and aggregated with:

```bash
python analysis/aggregate_conflicts.py --pretty experiments/results/*.jsonl
```

Runtime artifacts such as `leader.lock` and ad-hoc `events.log` files are intentionally not tracked.

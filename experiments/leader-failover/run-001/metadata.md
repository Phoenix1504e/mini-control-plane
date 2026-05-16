# Experiment Metadata

Date: 2026-05-16

Scenario:
Leader failover during delayed reconciliation

Configuration:
- Watch drop rate: 0.25
- Reconcile delay: 5s
- Controllers: 2
- etcd: single-node local
- Sampling interval: 250ms

Observed Behavior:
- Initial leader failure succeeded
- Standby promoted successfully
- Earlier resources remained unconverged
- demo3 converged successfully

Key Observations:
- Persistent drift under earlier failures
- Successful convergence after fixes
- MVCC semantics functioning correctly

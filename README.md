# Mini Control Plane

Mini Control Plane is an educational, Kubernetes-inspired control plane written in Go.
It demonstrates the core architectural patterns used by Kubernetes and other CNCF
projects, including declarative APIs, admission control, reconciliation loops,
watch-based informers, leader election, and status management.

This project is intentionally focused on *control-plane mechanics* rather than
workload execution or container orchestration.

---

## Motivation

Kubernetes is widely used, but its internal control-plane behavior is often treated
as a black box. Many engineers can operate Kubernetes without fully understanding
how its controllers, admission logic, and reconciliation model actually work.

Mini Control Plane exists to:
- make Kubernetes control-plane concepts explicit
- provide a readable reference implementation
- serve as a learning platform for cloud-native system design

The project favors clarity and correctness over scale and production optimizations.

---

## Design Goals

- Declarative, spec-driven system
- Clear separation between desired state (spec) and observed state (status)
- Idempotent reconciliation
- Event-driven controllers using watch & informers
- Policy enforcement via admission control
- Safe multi-controller operation through leader election
- Minimal dependencies and readable Go code

---

## High-Level Architecture

```text
Client
  |
  v
API Server
  |  (Admission)
  v
Storage
  |
  v
Watch / Informer
  |
  v
Controller
  |
  v
Runtime State
```



The API server is the single entry point for user intent.
Controllers react to state changes and converge the system toward the desired state.

---

## Core Components

### API Server
Acts as the control-plane front door.
Responsible for accepting resource definitions, enforcing admission policies,
persisting desired state, and emitting watch events.

### Admission Controller
Validates resources at creation time and records admission outcomes as status conditions.

### Storage Layer
Abstracted storage interface inspired by etcd semantics.
Current implementation uses file-backed storage.

### Watch & Informers
Controllers subscribe to resource change events instead of polling, enabling
immediate reconciliation.

### Controller
Implements reconciliation logic, drift detection, and status updates.

### Leader Election
Ensures only one controller instance actively reconciles at a time.
## Resource Model

Resources follow Kubernetes-style structure:

```yaml
spec:
  name: app-example
  replicas: 3
status:
  currentReplicas: 3
  conditions:
    - type: AdmissionApproved
      status: "True"
      reason: Accepted
```
- spec is user-owned and immutable by controllers
- status is system-owned and reflects observed state
- Controllers never mutate spec

## Key Concepts Demonstrated

- Declarative desired state
- Eventual consistency
- Idempotent reconciliation
- Drift detection and self-healing
- Admission-time policy enforcement
- Watch-based informers
- Leader election and safety

## Non-Goals

- Full Kubernetes API compatibility
- Production-grade scalability
- Container scheduling or networking
- etcd reimplementation

This project is intentionally scoped to control-plane behavior rather than
full orchestration.

---

## Project Status

**Stage:** Experimental / Educational  
**Maturity:** CNCF Sandbox–style  
**Stability:** Not production-ready

---

## Governance

This project follows lightweight governance inspired by CNCF principles.

### Maintainer
- Aditya Pathak (initial author and maintainer)

### Decision Making
- Technical decisions are made by the maintainer
- Design discussions are documented through issues

### Contributions
Contributions are welcome in the form of:
- design discussions
- documentation improvements
- bug reports
- feature proposals

This project follows the CNCF Code of Conduct:
https://github.com/cncf/foundation/blob/main/code-of-conduct.md

---

## Roadmap

### Short Term
- Resource versioning and optimistic concurrency
- Improved error handling and status reporting
- Structured event objects
- Expanded test coverage

### Medium Term
- Mutating admission controllers
- API versioning (v1alpha1, v1beta1)
- Simple scheduler abstraction
- Pluggable storage backends

### Long Term
- RBAC-style authorization
- Multi-resource dependency graph
- Controller crash recovery
- Educational visualizations of reconciliation

---

## License

Apache License 2.0

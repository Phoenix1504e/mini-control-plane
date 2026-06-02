# Paper

Mini Control Plane accompanies the paper:

> White-Box Fault Analysis for Kubernetes-Style Control Plane Semantics.

Add the final conference DOI, ACM/IEEE landing page, or arXiv/preprint URL here when it is available. The repository is organized so reviewers can connect the paper claims to code and artifacts:

- `pkg/fault/` contains the fault-injection middleware surface.
- `specs/experiments/` contains the Table 1 experiment configurations.
- `analysis/aggregate_conflicts.py` aggregates structured conflict logs.
- `experiments/results/` is reserved for JSONL outputs from controlled runs.

## Citation

```bibtex
@inproceedings{pathak2026minicontrolplane,
  title = {White-Box Fault Analysis for Kubernetes-Style Control Plane Semantics},
  author = {Pathak, Aditya},
  booktitle = {Proceedings of the International Conference on Distributed Computing and Networking},
  year = {2026},
  note = {Artifact: https://github.com/Phoenix1504e/mini-control-plane}
}
```

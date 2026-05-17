import json
import os
import matplotlib.pyplot as plt

BASE = "experiments/mvcc-conflicts"

controllers = []
conflicts = []

for run in sorted(os.listdir(BASE)):

    run_path = os.path.join(BASE, run)

    meta_path = os.path.join(run_path, "metadata.json")
    conflict_path = os.path.join(run_path, "mvcc_conflicts.jsonl")

    if not os.path.exists(meta_path):
        continue

    with open(meta_path) as f:
        meta = json.load(f)

    controller_count = meta["controllers"]

    conflict_count = 0

    if os.path.exists(conflict_path):
        with open(conflict_path) as f:
            for _ in f:
                conflict_count += 1

    controllers.append(controller_count)
    conflicts.append(conflict_count)

plt.figure(figsize=(8, 5))

plt.plot(
    controllers,
    conflicts,
    marker='o',
)

plt.xlabel("Controllers")
plt.ylabel("MVCC Conflicts")
plt.title("MVCC Conflict Scaling")

plt.grid(True)

plt.savefig("analysis/conflict_scaling.png")

print("Saved plot to analysis/conflict_scaling.png")

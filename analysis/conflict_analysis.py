import json
import os
from collections import defaultdict

BASE = "experiments/mvcc-conflicts"

results = []

for run in sorted(os.listdir(BASE)):

    run_path = os.path.join(BASE, run)

    meta_path = os.path.join(run_path, "metadata.json")
    conflict_path = os.path.join(run_path, "mvcc_conflicts.jsonl")

    if not os.path.exists(meta_path):
        continue

    with open(meta_path) as f:
        meta = json.load(f)

    controllers = meta["controllers"]

    conflicts = 0

    if os.path.exists(conflict_path):

        with open(conflict_path) as f:
            for _ in f:
                conflicts += 1

    results.append({
        "controllers": controllers,
        "conflicts": conflicts,
    })

print("\nMVCC Conflict Scaling Results\n")

for r in results:
    print(
        f"controllers={r['controllers']} "
        f"conflicts={r['conflicts']}"
    )

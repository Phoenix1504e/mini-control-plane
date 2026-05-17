import json
import math
import os
from collections import defaultdict

BASE = "experiments/mvcc-conflicts/campaign-001"

results = defaultdict(list)

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

    results[controllers].append(conflicts)

print("\nMVCC Conflict Aggregation\n")

for controllers in sorted(results.keys()):

    vals = results[controllers]

    avg = sum(vals) / len(vals)

    variance = sum((x - avg) ** 2 for x in vals) / len(vals)

    stddev = math.sqrt(variance)

    print(
        f"controllers={controllers} "
        f"trials={len(vals)} "
        f"avg={avg:.2f} "
        f"stddev={stddev:.2f} "
        f"runs={vals}"
    )

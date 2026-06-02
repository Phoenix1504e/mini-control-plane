#!/usr/bin/env python3
"""Aggregate structured conflict events emitted by experiment runs."""

from __future__ import annotations

import argparse
import json
from collections import Counter, defaultdict
from pathlib import Path
from typing import Any, Iterable


def load_events(paths: Iterable[Path]) -> Iterable[dict[str, Any]]:
    for path in paths:
        with path.open(encoding="utf-8") as handle:
            for line_number, line in enumerate(handle, 1):
                line = line.strip()
                if not line:
                    continue
                try:
                    yield json.loads(line)
                except json.JSONDecodeError as exc:
                    raise ValueError(f"{path}:{line_number}: invalid JSON: {exc}") from exc


def is_conflict(event: dict[str, Any]) -> bool:
    return event.get("event") == "conflict" or event.get("error") == "resourceVersion mismatch"


def aggregate(events: Iterable[dict[str, Any]]) -> dict[str, Any]:
    by_controller: Counter[str] = Counter()
    by_resource: Counter[str] = Counter()
    by_run: dict[str, int] = defaultdict(int)
    total = 0

    for event in events:
        if not is_conflict(event):
            continue
        total += 1
        by_controller[str(event.get("controller", "unknown"))] += 1
        by_resource[str(event.get("resource", "unknown"))] += 1
        by_run[str(event.get("run", "unknown"))] += 1

    return {
        "total_conflicts": total,
        "by_controller": dict(sorted(by_controller.items())),
        "by_resource": dict(sorted(by_resource.items())),
        "by_run": dict(sorted(by_run.items())),
    }


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("logs", nargs="+", type=Path, help="JSONL log files to aggregate")
    parser.add_argument("--pretty", action="store_true", help="Pretty-print the JSON summary")
    args = parser.parse_args()

    summary = aggregate(load_events(args.logs))
    indent = 2 if args.pretty else None
    print(json.dumps(summary, indent=indent, sort_keys=True))


if __name__ == "__main__":
    main()

import os
import json
import glob
import argparse
from datetime import datetime
import numpy as np

def parse_trace_journal(file_path):
    timestamps = []
    with open(file_path, 'r') as f:
        for line in f:
            if not line.strip(): continue
            try:
                event = json.loads(line)
                # Convert ISO timestamp string to datetime object
                ts_str = event.get("timestamp").replace("Z", "")
                dt = datetime.fromisoformat(ts_str)
                timestamps.append(dt)
            except (json.JSONDecodeError, ValueError):
                continue
    
    if len(timestamps) < 2:
        return 0
    
    # Calculate total duration of the trial in milliseconds
    duration = (max(timestamps) - min(timestamps)).total_seconds() * 1000
    return duration

def analyze_campaign(parent_dir):
    trace_files = sorted(glob.glob(os.path.join(parent_dir, "**/*.jsonl"), recursive=True))
    
    if not trace_files:
        print(f"[-] No trace journals found in {parent_dir}.")
        return

    durations = []
    print(f"| {'Run':<15} | {'Trial Duration (ms)':<20} |")
    print("-" * 40)

    for file_path in trace_files:
        dur = parse_trace_journal(file_path)
        if dur > 0:
            durations.append(dur)
            run_name = os.path.basename(os.path.dirname(file_path))
            print(f"| {run_name:<15} | {dur:<20.2f} |")

    if durations:
        print("=" * 40)
        print(f"Global Average Duration: {np.mean(durations):.2f} ms")
        print(f"Std Deviation          : {np.std(durations):.2f} ms")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--run-dir", required=True)
    args = parser.parse_args()
    analyze_campaign(args.run_dir)

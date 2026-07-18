import json
import numpy as np
import os
from datetime import datetime

def calculate_all_trials():
    all_iats = []
    
    # Iterate through all 10 runs
    for i in range(1, 11):
        path = f"experiments/fault-class-a/run-{i}/events.jsonl"
        if not os.path.exists(path):
            continue
            
        timestamps = []
        with open(path, 'r') as f:
            for line in f:
                event = json.loads(line)
                # Parse ISO string to datetime object
                dt = datetime.fromisoformat(event['timestamp'].replace('Z', '+00:00'))
                # Convert to total seconds (or milliseconds)
                timestamps.append(dt.timestamp() * 1000)
        
        # Calculate gaps (IATs) for this trial
        timestamps.sort()
        if len(timestamps) > 1:
            gaps = np.diff(timestamps)
            all_iats.extend(gaps)

    if not all_iats:
        print("No valid data found.")
        return

    # Calculate statistics
    print("--- Aggregated Performance Metrics (10 Trials) ---")
    print(f"Mean IAT     : {np.mean(all_iats):.4f} ms")
    print(f"Jitter (Std) : {np.std(all_iats):.4f} ms")
    print(f"P95 Latency  : {np.percentile(all_iats, 95):.4f} ms")
    print(f"P99 Latency  : {np.percentile(all_iats, 99):.4f} ms")
    print("--------------------------------------------------")

if __name__ == "__main__":
    calculate_all_trials()

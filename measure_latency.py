import json
import statistics
from datetime import datetime

def calculate_metrics(file_path):
    timestamps = []
    
    # Read the file and collect all success event timestamps
    with open(file_path, 'r') as f:
        for line in f:
            try:
                event = json.loads(line)
                if event.get("message") == "CreateSucceeded":
                    ts = datetime.fromisoformat(event["timestamp"].replace("Z", ""))
                    timestamps.append(ts)
            except (json.JSONDecodeError, ValueError, KeyError):
                continue
    
    # We need at least two events to calculate a gap
    if len(timestamps) < 2:
        return None
    
    # Sort timestamps to ensure they are in order
    timestamps.sort()
    
    # Calculate the gaps between each consecutive success in milliseconds
    gaps = [(timestamps[i] - timestamps[i-1]).total_seconds() * 1000 
            for i in range(1, len(timestamps))]
    
    return gaps

if __name__ == "__main__":
    # Point this to the log file of the trial you want to analyze
    path = "experiments/fault-class-a/run-1/events.jsonl"
    
    gaps = calculate_metrics(path)
    
    if gaps:
        avg_iat = sum(gaps) / len(gaps)
        jitter = statistics.stdev(gaps) if len(gaps) > 1 else 0
        
        print("--- Control Plane Performance Metrics ---")
        print(f"Total Successful Operations: {len(gaps) + 1}")
        print(f"Average Inter-arrival Time : {avg_iat:.2f} ms")
        print(f"Jitter (Std Deviation)     : {jitter:.2f} ms")
        print("------------------------------------------")
    else:
        print("[-] No 'CreateSucceeded' events found in the log.")

import os
import subprocess
import time
import shutil
import requests # Ensure you have 'requests' installed (pip install requests)

def run_campaign(total_trials=10, concurrency=5, ops=100):
    base_run_dir = "experiments/fault-class-a"
    os.makedirs(base_run_dir, exist_ok=True)

    for trial in range(1, total_trials + 1):
        # 1. Clean state
        if os.path.exists("default.etcd"):
            shutil.rmtree("default.etcd")
        
        # 2. Start etcd in the background
        etcd_proc = subprocess.Popen(
            ["etcd", "--listen-client-urls", "http://0.0.0.0:2379", 
             "--advertise-client-urls", "http://127.0.0.1:2379"],
            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL
        )
        
        # 3. Wait for etcd to be healthy
        print(f"[*] Starting Trial {trial}/{total_trials}...", end="", flush=True)
        ready = False
        for _ in range(10): # Try for 5 seconds
            try:
                if requests.get("http://localhost:2379/health").status_code == 200:
                    ready = True
                    break
            except:
                time.sleep(0.5)
        
        if not ready:
            print(" FAILED: etcd failed to start.")
            etcd_proc.terminate()
            break
            
        # 4. Run bench
        trial_dir = os.path.join(base_run_dir, f"run-{trial}")
        os.makedirs(trial_dir, exist_ok=True)
        output_path = os.path.join(trial_dir, "events.jsonl")
        
        cmd = ["./bin/bench", f"--concurrency={concurrency}", f"--ops={ops}", f"--out={output_path}"]
        subprocess.run(cmd, check=True, stdout=subprocess.DEVNULL)
        
        # 5. Clean shutdown
        etcd_proc.terminate()
        etcd_proc.wait()
        print(" Success.")

    print(f"\n[+] Campaign complete! Data saved in {base_run_dir}")

if __name__ == "__main__":
    run_campaign()

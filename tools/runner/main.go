package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	scenario := flag.String("scenario", "default", "Name of the experiment scenario")
	resources := flag.Int("resources", 10, "Number of simulated resources to create")
	flag.Parse()

	runID := fmt.Sprintf("run-%d", time.Now().Unix())
	archiveDir := filepath.Join("experiments", *scenario, runID)

	log.Printf("Starting experiment scenario: %s (Run ID: %s)", *scenario, runID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Start etcd daemon process
	etcdCmd := exec.CommandContext(ctx, "etcd",
		"--listen-client-urls", "http://127.0.0.1:2379",
		"--advertise-client-urls", "http://127.0.0.1:2379",
		"--log-level", "warn",
	)
	if err := etcdCmd.Start(); err != nil {
		log.Fatalf("Failed to spin up etcd: %v", err)
	}
	log.Println("✔ etcd instance deployed successfully")
	time.Sleep(2 * time.Second) // Grant etcd settling time

	// 2. Start Compiled Core API Server Binary
	apiCmd := exec.CommandContext(ctx, "./bin/apiserver")
	apiCmd.Stdout = os.Stdout
	apiCmd.Stderr = os.Stderr
	if err := apiCmd.Start(); err != nil {
		log.Fatalf("Failed to initialize API Server: %v", err)
	}
	log.Println("✔ API Server initialized on :8080")

	// 3. Start Compiled Controller Binary
	ctrlCmd := exec.CommandContext(ctx, "./bin/controller")
	ctrlCmd.Stdout = os.Stdout
	ctrlCmd.Stderr = os.Stderr
	if err := ctrlCmd.Start(); err != nil {
		log.Fatalf("Failed to initialize Controller loop: %v", err)
	}
	log.Println("✔ Control plane loops active")
	time.Sleep(5 * time.Second)

	// 4. Inject simulated resource workloads to stimulate high-conflict races
	log.Printf("Injecting %d resource creation requests into API server...", *resources)
	var wg sync.WaitGroup
	for i := 1; i <= *resources; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resName := fmt.Sprintf("res-%d", id)

			// FIX: Explicitly passing "replicas": 3 to stimulate the reconciliation and scheduling logic
			jsonPayload := fmt.Sprintf(`{"metadata":{"name":"%s"},"spec":{"name":"%s","replicas":3}}`, resName, resName)

			// Point to correct endpoint path and pass proper media tags
			_ = exec.Command("curl", "-s", "-X", "POST",
				"-H", "Content-Type: application/json",
				"-d", jsonPayload,
				"http://127.0.0.1:8080/resource").Run()
		}(i)
	}
	wg.Wait()
	log.Println("✔ Resource workload injections processed completely")

	// 5. Allow telemetry hooks to capture trailing background reconciliation loops
	log.Println("Holding for trailing reconciliation events to balance telemetry records...")
	time.Sleep(5 * time.Second)

	// 6. Graceful component winding down via signal cancellations
	log.Println("Signaling termination down to running control plane clusters...")
	cancel()

	// Await absolute platform cleanup resolution
	_ = apiCmd.Wait()
	_ = ctrlCmd.Wait()
	_ = etcdCmd.Wait()

	// Give the operating system flush buffers a tiny window to finish writing
	time.Sleep(2 * time.Second)

	// 7. Dynamic telemetry telemetry archiving segment
	log.Printf("Archiving telemetry telemetry outputs down to destination track: %s", archiveDir)
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		log.Fatalf("Failed to construct directory hierarchy structure: %v", err)
	}

	matches, _ := filepath.Glob("*.jsonl")
	if len(matches) == 0 {
		log.Println("⚠ Warning: No runtime log sequences (*.jsonl) matched in workspace root!")
		return
	}

	for _, file := range matches {
		dest := filepath.Join(archiveDir, file)
		if err := os.Rename(file, dest); err != nil {
			log.Printf("Failed to archive telemetry data target %s -> %s: %v", file, dest, err)
		} else {
			log.Printf("Successfully archived tracking segment: %s", file)
		}
	}
	log.Println("✔ Operational pipeline sequence finished cleanly.")
}

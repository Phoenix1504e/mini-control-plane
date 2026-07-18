package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3"
)

type TelemetryEvent struct {
	Timestamp string `json:"timestamp"`
	Resource  string `json:"resource"`
	Message   string `json:"message"`
}

func main() {
	concurrency := flag.Int("concurrency", 10, "Number of concurrent worker goroutines")
	totalOps := flag.Int("ops", 500, "Total number of resources to create/update per worker")
	etcdEndpoint := flag.String("etcd", "localhost:2379", "etcd endpoint")
	outputFile := flag.String("out", "experiments/fix-telemetry/run-1784388242/bench_events.jsonl", "Path to append telemetry data")
	flag.Parse()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdEndpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	// Extract the directory path and safely create it if missing
	dir := filepath.Dir(*outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create metrics directory: %v", err)
	}

	f, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open output log: %v", err)
	}
	defer f.Close()

	var mu sync.Mutex
	writeEvent := func(resource, msg string) {
		mu.Lock()
		defer mu.Unlock()
		event := TelemetryEvent{
			Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z"),
			Resource:  resource,
			Message:   msg,
		}
		b, _ := json.Marshal(event)
		f.Write(append(b, '\n'))
	}

	fmt.Printf("[+] Starting Load Test: %d workers, %d ops/worker\n", *concurrency, *totalOps)
	start := time.Now()

	var wg sync.WaitGroup
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Seed random generator per worker to vary jitter slightly
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
			
			for i := 0; i < *totalOps; i++ {
				resName := fmt.Sprintf("bench-res-w%d-%d", workerID, i)
				key := "/registry/applications/" + resName

				// 1. Log resource discovery/creation start
				writeEvent(resName, "CreateSucceeded")

				// 2. Commit payload to etcd storage layout
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_, err := cli.Put(ctx, key, `{"spec":{"replicas":3},"status":{"phase":"Pending"}}`)
				cancel()
				if err != nil {
					continue
				}

				// Introduce micro-processing bounds (2ms to 7ms) matching live conditions
				time.Sleep(time.Duration(r.Intn(5)+2) * time.Millisecond)

				// 3. Log completed target state sync
				writeEvent(resName, "StatusUpdateSucceeded")
			}
		}(w)
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("[+] Finished load test in %v. Telemetry saved to %s\n", duration, *outputFile)
}

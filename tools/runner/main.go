package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Metadata struct {
	Scenario    string    `json:"scenario"`
	Controllers int       `json:"controllers"`
	Resources   int       `json:"resources"`
	Duration    string    `json:"duration"`
	StartedAt   time.Time `json:"startedAt"`
}

func main() {

	scenario := flag.String(
		"scenario",
		"mvcc-conflicts",
		"experiment scenario",
	)

	controllers := flag.Int(
		"controllers",
		4,
		"number of controllers",
	)

	resources := flag.Int(
		"resources",
		20,
		"number of resources",
	)

	duration := flag.Duration(
		"duration",
		30*time.Second,
		"experiment duration",
	)

	flag.Parse()

	runID := fmt.Sprintf(
		"run-%d",
		time.Now().Unix(),
	)

	runDir := filepath.Join(
		"experiments",
		*scenario,
		runID,
	)

	if err := os.MkdirAll(runDir, 0755); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting experiment", runID)

	// Start etcd
	etcd := exec.Command("etcd")
	etcd.Stdout = os.Stdout
	etcd.Stderr = os.Stderr

	if err := etcd.Start(); err != nil {
		log.Fatal(err)
	}

	defer etcd.Process.Kill()

	log.Println("Started etcd")

	time.Sleep(3 * time.Second)

	// Start controllers
	var controllerProcs []*exec.Cmd

	for i := 0; i < *controllers; i++ {

		cmd := exec.Command(
			"go",
			"run",
			"./cmd/controller",
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		controllerProcs = append(
			controllerProcs,
			cmd,
		)
	}

	log.Printf(
		"Started %d controllers",
		*controllers,
	)

	time.Sleep(5 * time.Second)

	// Generate workload
	for i := 0; i < *resources; i++ {

		name := fmt.Sprintf(
			"resource-%d",
			i,
		)

		payload := fmt.Sprintf(
			`{"metadata":{"name":"%s"},"spec":{"replicas":3}}`,
			name,
		)

		put := exec.Command(
			"etcdctl",
			"put",
			"/resources/"+name,
			payload,
		)

		if out, err := put.CombinedOutput(); err != nil {
			log.Printf(
				"workload error: %v %s",
				err,
				string(out),
			)
		}
	}

	log.Println("Workload generation complete")

	// Wait for experiment duration
	time.Sleep(*duration)

	// Stop controllers
	for _, proc := range controllerProcs {
		_ = proc.Process.Kill()
	}

	log.Println("Controllers stopped")

	// Archive telemetry
	copyFile("reconcile.jsonl", filepath.Join(runDir, "reconcile.jsonl"))
	copyFile("state_samples.jsonl", filepath.Join(runDir, "state_samples.jsonl"))
	copyFile("mvcc_conflicts.jsonl", filepath.Join(runDir, "mvcc_conflicts.jsonl"))
	copyFile("leader_events.jsonl", filepath.Join(runDir, "leader_events.jsonl"))

	meta := Metadata{
		Scenario:    *scenario,
		Controllers: *controllers,
		Resources:   *resources,
		Duration:    duration.String(),
		StartedAt:   time.Now(),
	}

	f, err := os.Create(
		filepath.Join(runDir, "metadata.json"),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	if err := enc.Encode(meta); err != nil {
		log.Fatal(err)
	}

	log.Println("Experiment archived at", runDir)
}

func copyFile(src, dst string) {

	in, err := os.ReadFile(src)
	if err != nil {
		return
	}

	_ = os.WriteFile(dst, in, 0644)
}

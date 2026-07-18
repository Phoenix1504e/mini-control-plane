package reconciler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.etcd.io/etcd/client/v3"
	"k8s.io/client-go/util/workqueue"
)

// Reconciler handles the control loop logic
type Reconciler struct {
	client *clientv3.Client
	queue  workqueue.RateLimitingInterface
}

func New(client *clientv3.Client, queue workqueue.RateLimitingInterface) *Reconciler {
	return &Reconciler{
		client: client,
		queue:  queue,
	}
}

// Reconcile processes a single item from the workqueue
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	log.Printf("[*] Reconciling key: %s", key)

	// 1. Fetch current state from etcd
	resp, err := r.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to fetch state: %w", err)
	}

	// 2. Idempotent logic: Drive observed state toward desired state
	// (Example placeholder: add your resource logic here)
	desiredState := "active"
	
	// 3. Perform atomic update with retry logic (Task 3 integration)
	err = r.updateWithBackoff(ctx, key, desiredState, 5)
	if err != nil {
		return fmt.Errorf("reconciliation failed after retries: %w", err)
	}

	log.Printf("[+] Reconciliation success for: %s", key)
	return nil
}

// updateWithBackoff implements exponential backoff to mitigate MVCC contention
func (r *Reconciler) updateWithBackoff(ctx context.Context, key, val string, maxAttempts int) error {
	baseDelay := 50 * time.Millisecond

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Attempting atomic Put or Txn
		_, err := r.client.Put(ctx, key, val)

		if err == nil {
			return nil // Success
		}

		// Log conflict for instrumentation metrics
		log.Printf("[!] Conflict detected on key %s (attempt %d/%d)", key, attempt+1, maxAttempts)

		// Calculate backoff: (base * 2^attempt) + jitter
		jitter := time.Duration(rand.Intn(50)) * time.Millisecond
		delay := (baseDelay * time.Duration(1<<attempt)) + jitter

		select {
		case <-time.After(delay):
			continue // Retry
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("exhausted %d retries due to storage conflicts", maxAttempts)
}

// Run starts the reconciliation loop worker
func (r *Reconciler) Run(ctx context.Context) {
	for {
		item, quit := r.queue.Get()
		if quit {
			return
		}

		key := item.(string)
		if err := r.Reconcile(ctx, key); err != nil {
			log.Printf("[-] Error reconciling %s: %v", key, err)
			r.queue.AddRateLimited(key)
		} else {
			r.queue.Forget(item)
		}
		r.queue.Done(item)
	}
}

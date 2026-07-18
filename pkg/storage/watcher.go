package storage

import (
	"context"
	"log"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// WatchWithFaults initiates a watch stream on the root prefix with a probabilistic drop rate.
// dropRate should be a float between 0.0 (no drops) and 1.0 (all drops).
func (s *EtcdStorage) WatchWithFaults(ctx context.Context, dropRate float64) clientv3.WatchChan {
	// Initialize the standard etcd watch channel
	watchChan := s.cli.Watch(ctx, s.root+"/", clientv3.WithPrefix())

	// Create a new channel to relay filtered events
	faultyChan := make(clientv3.WatchChan)

	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		defer close(faultyChan)
		for {
			select {
			case <-ctx.Done():
				return
			case resp, ok := <-watchChan:
				if !ok {
					return
				}

				// Fault Injection: Simulate "Watch Event Loss"
				if r.Float64() < dropRate {
					log.Printf("[!] Fault Injection: Dropped watch event at revision %d", resp.Header.Revision)
					continue
				}

				// Relay the event if it passes the fault gate
				faultyChan <- resp
			}
		}
	}()

	return faultyChan
}

// Watch provides the default clean watch stream (no faults)
func (s *EtcdStorage) Watch(ctx context.Context) clientv3.WatchChan {
	return s.WatchWithFaults(ctx, 0.0)
}

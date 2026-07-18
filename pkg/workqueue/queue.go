package workqueue

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type Queue struct {
	ch chan string

	retries map[string]int

	// keys waiting in queue
	dirty map[string]bool

	// keys currently being processed
	processing map[string]bool

	// Global strategy configuration ("baseline" or "jitter-backoff")
	strategy string

	mu sync.Mutex
}

func New(size int, strategy string) *Queue {
	return &Queue{
		ch:         make(chan string, size),
		retries:    make(map[string]int),
		dirty:      make(map[string]bool),
		processing: make(map[string]bool),
		strategy:   strategy,
	}
}

func (q *Queue) Add(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Already queued
	if q.dirty[key] {
		return
	}

	// Already being processed
	if q.processing[key] {
		return
	}

	q.dirty[key] = true
	q.ch <- key
}

func (q *Queue) AddRateLimited(key string) {
	q.mu.Lock()

	retry := q.retries[key]
	q.retries[key] = retry + 1

	delete(q.processing, key)
	strategy := q.strategy

	q.mu.Unlock()

	var delay time.Duration

	switch strategy {
	case "jitter-backoff":
		// Strategy 1: Exponential Backoff with Full Jitter
		baseDelay := 10.0      // 10ms base
		maxDelay := 1000.0     // 1000ms ceiling
		
		// Calculate exponential step: base * 2^retry
		temp := baseDelay * math.Pow(2, float64(retry))
		if temp > maxDelay {
			temp = maxDelay
		}
		
		// Full Jitter: Randomize between [0, temp) to split up colliding reconcilers
		seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
		jitteredMs := seededRand.Float64() * temp
		delay = time.Duration(jitteredMs) * time.Millisecond

	default:
		// Strategy 0: Baseline (Immediate Requeue / Minimal constant delay for serialization)
		delay = 0 * time.Millisecond
	}

	if delay == 0 {
		q.Add(key)
		return
	}

	go func() {
		time.Sleep(delay)
		q.Add(key)
	}()
}

func (q *Queue) Get() string {
	key := <-q.ch

	q.mu.Lock()

	delete(q.dirty, key)
	q.processing[key] = true

	q.mu.Unlock()

	return key
}

func (q *Queue) Forget(key string) {
	q.mu.Lock()

	delete(q.retries, key)
	delete(q.processing, key)

	q.mu.Unlock()
}

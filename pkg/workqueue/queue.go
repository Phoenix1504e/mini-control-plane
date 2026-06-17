package workqueue

import (
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

	mu sync.Mutex
}

func New(size int) *Queue {
	return &Queue{
		ch:         make(chan string, size),
		retries:    make(map[string]int),
		dirty:      make(map[string]bool),
		processing: make(map[string]bool),
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

	q.mu.Unlock()

	delay := time.Duration(retry+1) * time.Second

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

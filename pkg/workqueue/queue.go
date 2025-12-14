package workqueue

import (
	"sync"
	"time"
)

type Queue struct {
	ch      chan string
	retries map[string]int
	mu      sync.Mutex
}

func New(size int) *Queue {
	return &Queue{
		ch:      make(chan string, size),
		retries: make(map[string]int),
	}
}

func (q *Queue) Add(key string) {
	q.ch <- key
}

func (q *Queue) AddRateLimited(key string) {
	q.mu.Lock()
	retry := q.retries[key]
	q.retries[key] = retry + 1
	q.mu.Unlock()

	delay := time.Duration(retry+1) * time.Second

	go func() {
		time.Sleep(delay)
		q.ch <- key
	}()
}

func (q *Queue) Get() string {
	return <-q.ch
}

func (q *Queue) Forget(key string) {
	q.mu.Lock()
	delete(q.retries, key)
	q.mu.Unlock()
}

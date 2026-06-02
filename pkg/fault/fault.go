package fault

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Operation string

const (
	OperationCreate       Operation = "create"
	OperationUpdate       Operation = "update"
	OperationUpdateStatus Operation = "update_status"
	OperationGet          Operation = "get"
	OperationList         Operation = "list"
)

type Injection struct {
	Operation Operation
	Delay     time.Duration
	Error     error
	Rate      float64
}

type Middleware struct {
	mu         sync.RWMutex
	injections map[Operation]Injection
	rand       *rand.Rand
}

func New(seed int64) *Middleware {
	return &Middleware{
		injections: make(map[Operation]Injection),
		rand:       rand.New(rand.NewSource(seed)),
	}
}

func (m *Middleware) Configure(injection Injection) error {
	if injection.Rate < 0 || injection.Rate > 1 {
		return errors.New("fault injection rate must be between 0 and 1")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.injections[injection.Operation] = injection
	return nil
}

func (m *Middleware) Clear(operation Operation) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.injections, operation)
}

func (m *Middleware) Apply(ctx context.Context, operation Operation, next func(context.Context) error) error {
	injection, ok := m.pick(operation)
	if ok {
		if injection.Delay > 0 {
			timer := time.NewTimer(injection.Delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
			}
		}

		if injection.Error != nil {
			return injection.Error
		}
	}

	return next(ctx)
}

func (m *Middleware) pick(operation Operation) (Injection, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	injection, ok := m.injections[operation]
	if !ok {
		return Injection{}, false
	}

	if injection.Rate == 0 {
		return Injection{}, false
	}
	if injection.Rate < 1 && m.rand.Float64() > injection.Rate {
		return Injection{}, false
	}

	return injection, true
}

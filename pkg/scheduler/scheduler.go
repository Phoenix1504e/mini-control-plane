package scheduler

import (

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	"github.com/Phoenix1504e/mini-control-plane/pkg/storage"
)

type Scheduler struct {
	store storage.Storage
	nodes []string
}

func New(store storage.Storage, nodes []string) *Scheduler {
	return &Scheduler{
		store: store,
		nodes: nodes,
	}
}

// Schedule assigns replicas to nodes
func (s *Scheduler) Schedule(res *api.Resource) error {
	current, err := s.store.Get(res.Metadata.Name)
	if err != nil {
		return err
	}

	if current.Status.Placements == nil {
		current.Status.Placements = make(map[string]int)
	}

	desired := current.Spec.Replicas
	currentCount := current.Status.CurrentReplicas

	toSchedule := desired - currentCount
	if toSchedule <= 0 {
		return nil
	}

	i := 0
	for toSchedule > 0 {
		node := s.nodes[i%len(s.nodes)]
		current.Status.Placements[node]++
		current.Status.CurrentReplicas++
		toSchedule--
		i++
	}

	return s.store.UpdateStatus(current)
}

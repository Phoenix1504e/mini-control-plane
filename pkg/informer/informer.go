package informer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
)

type Event struct {
	Type     EventType
	Resource *api.Resource
}

type Informer struct {
	cli       *clientv3.Client
	root      string
	EventChan chan Event
}

func NewInformer(cli *clientv3.Client, root string) *Informer {
	return &Informer{
		cli:       cli,
		root:      root,
		EventChan: make(chan Event, 100),
	}
}

func (i *Informer) Start(ctx context.Context) {
	watchChan := i.cli.Watch(ctx, i.root+"/", clientv3.WithPrefix())

	for w := range watchChan {
		for _, ev := range w.Events {
			var res api.Resource

			switch ev.Type {
			case clientv3.EventTypePut:
				if err := json.Unmarshal(ev.Kv.Value, &res); err == nil {
					eventType := Modified
					if ev.Kv.CreateRevision == ev.Kv.ModRevision {
						eventType = Added
					}
					i.EventChan <- Event{Type: eventType, Resource: &res}
				}

			case clientv3.EventTypeDelete:
				if err := json.Unmarshal(ev.PrevKv.Value, &res); err == nil {
					i.EventChan <- Event{Type: Deleted, Resource: &res}
				}
			}
		}
	}

	log.Println("Informer stopped")
}

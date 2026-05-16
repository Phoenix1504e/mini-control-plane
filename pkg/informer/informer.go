package informer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Phoenix1504e/mini-control-plane/internal/faults"
	"github.com/Phoenix1504e/mini-control-plane/pkg/api"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EventType string

const (
	Added   EventType = "ADDED"
	Updated EventType = "UPDATED"
)

type Event struct {
	Type     EventType
	Resource *api.Resource
}

type Informer struct {
	client    *clientv3.Client
	prefix    string
	EventChan chan Event
}

func NewInformer(client *clientv3.Client, prefix string) *Informer {
	return &Informer{
		client:    client,
		prefix:    prefix,
		EventChan: make(chan Event),
	}
}

func (i *Informer) Start(ctx context.Context) {
	rawWatchChan := i.client.Watch(
		ctx,
		i.prefix,
		clientv3.WithPrefix(),
	)

	injector := &faults.ProbabilisticInjector{
		DropRate: 0.25,
	}

	watchChan := faults.WrapWatchChannel(
		ctx,
		rawWatchChan,
		injector,
	)

	for {
		select {

		case <-ctx.Done():
			return

		case resp, ok := <-watchChan:
			if !ok {
				log.Println("watch channel closed")
				return
			}

			for _, ev := range resp.Events {

				// Ignore invalid events
				if ev.Kv == nil || len(ev.Kv.Value) == 0 {
					continue
				}

				var res api.Resource

				if err := json.Unmarshal(ev.Kv.Value, &res); err != nil {
					log.Println("failed to decode resource:", err)
					continue
				}

				eventType := Updated

				if ev.Type == clientv3.EventTypePut && ev.IsCreate() {
					eventType = Added
				}

				i.EventChan <- Event{
					Type:     eventType,
					Resource: &res,
				}
			}
		}
	}
}


package informer

import (
	"context"
	"encoding/json"
	"log"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

type EventType string

const (
	Added   EventType = "ADDED"
	Updated EventType = "UPDATED"
	Deleted EventType = "DELETED"
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

	// ----------------------------------
	// Initial LIST
	// ----------------------------------
	resp, err := i.client.Get(
		ctx,
		i.prefix,
		clientv3.WithPrefix(),
	)
	if err != nil {
		log.Println("initial list failed:", err)
	} else {
		for _, kv := range resp.Kvs {
			var res api.Resource

			if err := json.Unmarshal(kv.Value, &res); err != nil {
				log.Println("failed to decode resource:", err)
				continue
			}

			i.EventChan <- Event{
				Type:     Added,
				Resource: &res,
			}
		}
	}

	// ----------------------------------
	// Watch for future updates
	// ----------------------------------
	watchChan := i.client.Watch(
		ctx,
		i.prefix,
		clientv3.WithPrefix(),
		clientv3.WithPrevKV(),
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

				var (
					res       api.Resource
					eventType EventType
					data      []byte
				)

				switch ev.Type {

				case clientv3.EventTypePut:
					data = ev.Kv.Value

					if ev.IsCreate() {
						eventType = Added
					} else {
						eventType = Updated
					}

				case clientv3.EventTypeDelete:
					if ev.PrevKv == nil {
						continue
					}

					data = ev.PrevKv.Value
					eventType = Deleted

				default:
					continue
				}

				if err := json.Unmarshal(data, &res); err != nil {
					log.Println("failed to decode resource:", err)
					continue
				}

				i.EventChan <- Event{
					Type:     eventType,
					Resource: &res,
				}
			}
		}
	}
}

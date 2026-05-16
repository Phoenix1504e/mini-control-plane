package leader

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/internal/telemetry"
)

type LeaderElector struct {
	client   *clientv3.Client
	key      string
	id       string
	isLeader bool

	logger *telemetry.Logger
}

func New(
	client *clientv3.Client,
	key,
	id string,
) *LeaderElector {

	logger, err := telemetry.NewLogger(
		"leader_events.jsonl",
	)

	if err != nil {
		log.Fatalf(
			"failed to initialize leader telemetry: %v",
			err,
		)
	}

	return &LeaderElector{
		client: client,
		key:    key,
		id:     id,
		logger: logger,
	}
}

// Run blocks until leadership is acquired, then maintains it
func (l *LeaderElector) Run(
	ctx context.Context,
) error {

	for {

		select {

		case <-ctx.Done():
			return nil

		default:
		}

		// Already leader
		if l.isLeader {

			time.Sleep(2 * time.Second)

			continue
		}

		// Try to acquire leadership
		leaseID, err := l.acquire(ctx)

		if err != nil {

			log.Println(
				"leader election error:",
				err,
			)

			time.Sleep(2 * time.Second)

			continue
		}

		if leaseID == 0 {

			// Another leader exists
			time.Sleep(2 * time.Second)

			continue
		}

		// Became leader
		log.Println("Leader elected:", l.id)

		l.isLeader = true

		err = l.logger.Emit(telemetry.Event{
			Timestamp: time.Now(),
			Type:      "leader_elected",
			LeaderID:  l.id,
		})

		if err != nil {
			log.Printf(
				"failed to emit leader_elected telemetry: %v",
				err,
			)
		}

		// Hold leadership
		l.holdLeadership(
			ctx,
			leaseID,
		)

		log.Println(
			"Leadership lost, entering standby mode",
		)

		l.isLeader = false

		err = l.logger.Emit(telemetry.Event{
			Timestamp: time.Now(),
			Type:      "leader_lost",
			LeaderID:  l.id,
		})

		if err != nil {
			log.Printf(
				"failed to emit leader_lost telemetry: %v",
				err,
			)
		}
	}
}

func (l *LeaderElector) acquire(
	ctx context.Context,
) (clientv3.LeaseID, error) {

	lease, err := l.client.Grant(ctx, 5)

	if err != nil {
		return 0, err
	}

	txn := l.client.Txn(ctx).
		If(
			clientv3.Compare(
				clientv3.CreateRevision(l.key),
				"=",
				0,
			),
		).
		Then(
			clientv3.OpPut(
				l.key,
				l.id,
				clientv3.WithLease(lease.ID),
			),
		)

	resp, err := txn.Commit()

	if err != nil {
		return 0, err
	}

	if !resp.Succeeded {
		return 0, nil
	}

	return lease.ID, nil
}

func (l *LeaderElector) holdLeadership(
	ctx context.Context,
	leaseID clientv3.LeaseID,
) {

	ch, err := l.client.KeepAlive(
		ctx,
		leaseID,
	)

	if err != nil {
		return
	}

	for {

		select {

		case <-ctx.Done():
			return

		case _, ok := <-ch:

			if !ok {

				// Lease expired
				return
			}
		}
	}
}

func (l *LeaderElector) IsLeader() bool {
	return l.isLeader
}

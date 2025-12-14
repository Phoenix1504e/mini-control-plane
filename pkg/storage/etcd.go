package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

type EtcdStorage struct {
	cli  *clientv3.Client
	root string
}

// NewEtcdStorage initializes an etcd-backed storage
func NewEtcdStorage(root string) (*EtcdStorage, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdStorage{
		cli:  cli,
		root: root,
	}, nil
}

// key returns the etcd key for a resource
func (s *EtcdStorage) key(name string) string {
	return fmt.Sprintf("%s/%s", s.root, name)
}

// Create stores a resource only if it does not already exist
func (s *EtcdStorage) Create(res *api.Resource) error {
	key := s.key(res.Spec.Name)

	val, err := json.Marshal(res)
	if err != nil {
		return err
	}

	txnResp, err := s.cli.Txn(context.Background()).
		If(clientv3.Compare(clientv3.Version(key), "=", 0)).
		Then(clientv3.OpPut(key, string(val))).
		Commit()

	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return fmt.Errorf("resource already exists")
	}

	return nil
}

// Update updates a resource using optimistic concurrency (resourceVersion)
func (s *EtcdStorage) Update(res *api.Resource) error {
	key := s.key(res.Spec.Name)

	if res.Metadata.ResourceVersion == "" {
		return fmt.Errorf("missing resourceVersion")
	}

	rv, err := strconv.ParseInt(res.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid resourceVersion")
	}

	val, err := json.Marshal(res)
	if err != nil {
		return err
	}

	txnResp, err := s.cli.Txn(context.Background()).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", rv)).
		Then(clientv3.OpPut(key, string(val))).
		Commit()

	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return fmt.Errorf("conflict: resourceVersion mismatch")
	}

	return nil
}

// Get retrieves a single resource and sets its resourceVersion
func (s *EtcdStorage) Get(name string) (*api.Resource, error) {
	key := s.key(name)

	resp, err := s.cli.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("not found")
	}

	kv := resp.Kvs[0]

	var res api.Resource
	if err := json.Unmarshal(kv.Value, &res); err != nil {
		return nil, err
	}

	res.Metadata.ResourceVersion = fmt.Sprintf("%d", kv.ModRevision)
	return &res, nil
}

// List retrieves all resources and populates resourceVersion
func (s *EtcdStorage) List() ([]*api.Resource, error) {
	resp, err := s.cli.Get(
		context.Background(),
		s.root+"/",
		clientv3.WithPrefix(),
	)
	if err != nil {
		return nil, err
	}

	var items []*api.Resource

	for _, kv := range resp.Kvs {
		var res api.Resource
		if err := json.Unmarshal(kv.Value, &res); err == nil {
			res.Metadata.ResourceVersion = fmt.Sprintf("%d", kv.ModRevision)
			items = append(items, &res)
		}
	}

	return items, nil
}
func (s *EtcdStorage) UpdateStatus(res *api.Resource) error {
	key := s.key(res.Spec.Name)

	if res.Metadata.ResourceVersion == "" {
		return fmt.Errorf("missing resourceVersion")
	}

	rv, err := strconv.ParseInt(res.Metadata.ResourceVersion, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid resourceVersion")
	}

	// Get latest object
	resp, err := s.cli.Get(context.Background(), key)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return fmt.Errorf("not found")
	}

	var current api.Resource
	if err := json.Unmarshal(resp.Kvs[0].Value, &current); err != nil {
		return err
	}

	// Only mutate STATUS
	current.Status = res.Status

	val, _ := json.Marshal(&current)

	txnResp, err := s.cli.Txn(context.Background()).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", rv)).
		Then(clientv3.OpPut(key, string(val))).
		Commit()

	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return fmt.Errorf("conflict updating status")
	}

	return nil
}

package storage

import "github.com/Phoenix1504e/mini-control-plane/pkg/api"

type Storage interface {
	Create(res *api.Resource) error
	Update(res *api.Resource) error
	UpdateStatus(res *api.Resource) error
	Delete(name string) error

	Get(name string) (*api.Resource, error)
	List() ([]*api.Resource, error)
}

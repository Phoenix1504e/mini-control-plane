package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

type FileStorage struct {
	BaseDir string
	WatchCh chan api.WatchEvent
}

func NewFileStorage(dir string) *FileStorage {
	_ = os.MkdirAll(dir, 0755)
	return &FileStorage{
		BaseDir: dir,
		WatchCh: make(chan api.WatchEvent, 100),
	}
}

func (fs *FileStorage) path(name string) string {
	return filepath.Join(fs.BaseDir, name+".yaml")
}

func (fs *FileStorage) Create(res *api.Resource) error {
	path := fs.path(res.Spec.Name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("resource already exists")
	}

	if err := fs.Update(res); err != nil {
		return err
	}

	fs.WatchCh <- api.WatchEvent{
		Type:     api.Added,
		Resource: res,
	}

	return nil
}

func (fs *FileStorage) Update(res *api.Resource) error {
	data, err := yaml.Marshal(res)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fs.path(res.Spec.Name), data, 0644); err != nil {
		return err
	}

	fs.WatchCh <- api.WatchEvent{
		Type:     api.Updated,
		Resource: res,
	}

	return nil
}

func (fs *FileStorage) Get(name string) (*api.Resource, error) {
	data, err := os.ReadFile(fs.path(name))
	if err != nil {
		return nil, err
	}

	var res api.Resource
	if err := yaml.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (fs *FileStorage) List() ([]*api.Resource, error) {
	files, err := os.ReadDir(fs.BaseDir)
	if err != nil {
		return nil, err
	}

	var out []*api.Resource
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(fs.BaseDir, f.Name()))
		if err != nil {
			continue
		}

		var res api.Resource
		if yaml.Unmarshal(data, &res) == nil {
			out = append(out, &res)
		}
	}

	return out, nil
}

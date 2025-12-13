package store

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

func LoadResources(dir string) ([]*api.Resource, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var resources []*api.Resource

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var res api.Resource
		if err := yaml.Unmarshal(data, &res); err != nil {
			continue
		}

		resources = append(resources, &res)
	}

	return resources, nil
}

func SaveResource(path string, res *api.Resource) error {
	data, err := yaml.Marshal(res)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

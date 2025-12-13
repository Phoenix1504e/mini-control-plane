package runtime

import (
	"fmt"
	"os"
	"strings"
)

func ListInstances(name string) ([]string, error) {
	files, err := os.ReadDir("./runtime")
	if err != nil {
		return nil, err
	}

	var instances []string
	prefix := name + "-"

	for _, f := range files {
		if strings.HasPrefix(f.Name(), prefix) {
			instances = append(instances, f.Name())
		}
	}

	return instances, nil
}

func CreateInstance(name string, id int) error {
	return os.WriteFile(
		fmt.Sprintf("./runtime/%s-%d", name, id),
		[]byte("running"),
		0644,
	)
}

func DeleteInstance(name string) error {
	return os.Remove("./runtime/" + name)
}

func CountInstances(name string) (int, error) {
	instances, err := ListInstances(name)
	if err != nil {
		return 0, err
	}
	return len(instances), nil
}

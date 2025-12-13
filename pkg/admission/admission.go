package admission

import (
	"fmt"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

const maxReplicas = 5

func Validate(res *api.Resource) error {
	spec := res.Spec

	if spec.Replicas < 0 {
		return fmt.Errorf("replicas cannot be negative")
	}

	if spec.Replicas > maxReplicas {
		return fmt.Errorf("replicas exceed maximum allowed (%d)", maxReplicas)
	}

	if len(spec.Name) < 4 || spec.Name[:4] != "app-" {
		return fmt.Errorf("resource name must start with 'app-'")
	}

	return nil
}

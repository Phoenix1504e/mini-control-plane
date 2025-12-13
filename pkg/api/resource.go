package api

type Condition struct {
	Type               string `yaml:"type"`
	Status             string `yaml:"status"`
	Reason             string `yaml:"reason"`
	Message            string `yaml:"message"`
	LastTransitionTime string `yaml:"lastTransitionTime"`
}

type ResourceSpec struct {
	Name     string `yaml:"name"`
	Replicas int    `yaml:"replicas"`
}

type ResourceStatus struct {
	CurrentReplicas int         `yaml:"currentReplicas,omitempty"`
	LastReconciled  string      `yaml:"lastReconciled,omitempty"`
	Conditions      []Condition `yaml:"conditions,omitempty"`
}

type Resource struct {
	Spec   ResourceSpec   `yaml:"spec"`
	Status ResourceStatus `yaml:"status,omitempty"`
}

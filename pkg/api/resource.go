package api

type ObjectMeta struct {
	Name            string `json:"name"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
}

type ResourceSpec struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
}

type ResourceStatus struct {
	CurrentReplicas int         `json:"currentReplicas,omitempty"`
        Placements map[string]int `json:"placements,omitempty"`
	Conditions      []Condition `json:"conditions,omitempty"`
}

type Resource struct {
	Metadata ObjectMeta     `json:"metadata"`
	Spec     ResourceSpec   `json:"spec"`
	Status   ResourceStatus `json:"status,omitempty"`
}

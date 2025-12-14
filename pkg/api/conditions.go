package api

import "time"

// Condition represents the state of the resource at a certain point
type Condition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"` // True, False, Unknown
	Reason             string    `json:"reason,omitempty"`
	Message            string    `json:"message,omitempty"`
	LastTransitionTime time.Time `json:"lastTransitionTime,omitempty"`
}

// SetCondition adds or updates a condition in the list
func SetCondition(conditions *[]Condition, newCond Condition) {
	for i, cond := range *conditions {
		if cond.Type == newCond.Type {
			newCond.LastTransitionTime = time.Now()
			(*conditions)[i] = newCond
			return
		}
	}

	newCond.LastTransitionTime = time.Now()
	*conditions = append(*conditions, newCond)
}

// GetCondition returns a condition by type
func GetCondition(conditions []Condition, condType string) *Condition {
	for _, cond := range conditions {
		if cond.Type == condType {
			return &cond
		}
	}
	return nil
}

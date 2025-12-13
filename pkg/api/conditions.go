package api

import "time"

func SetCondition(
	conditions []Condition,
	newCond Condition,
) []Condition {
	for i, c := range conditions {
		if c.Type == newCond.Type {
			if c.Status != newCond.Status {
				newCond.LastTransitionTime = time.Now().Format(time.RFC3339)
			} else {
				newCond.LastTransitionTime = c.LastTransitionTime
			}
			conditions[i] = newCond
			return conditions
		}
	}

	newCond.LastTransitionTime = time.Now().Format(time.RFC3339)
	return append(conditions, newCond)
}

func IsConditionTrue(conditions []Condition, condType string) bool {
	for _, c := range conditions {
		if c.Type == condType && c.Status == "True" {
			return true
		}
	}
	return false
}

package api

type WatchEventType string

const (
	Added   WatchEventType = "ADDED"
	Updated WatchEventType = "UPDATED"
)

type WatchEvent struct {
	Type     WatchEventType
	Resource *Resource
}

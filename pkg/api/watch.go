package api

type WatchEventType string

const (
	Added   WatchEventType = "ADDED"
	Updated WatchEventType = "UPDATED"
	Deleted WatchEventType = "DELETED"
)

type WatchEvent struct {
	Type     WatchEventType
	Resource *Resource
}

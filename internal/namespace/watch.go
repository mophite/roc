package namespace

type WatcherAction string

const (
	WatcherCreate WatcherAction = "create"
	WatcherUpdate               = "update"
	WatcherDelete               = "delete"
)

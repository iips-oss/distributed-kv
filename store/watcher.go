package store

type event struct {
	command string
	key string
	value string
	revision int64
}

type watcher struct {
	prefix string
	events chan event
}
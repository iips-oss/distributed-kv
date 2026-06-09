package store

// Using iota to define commands, as strings typos would compile silently
type Command int

const (
	SET Command = iota
	// GET
	DELETE
)

// Event has to be exported because it would be used outside of store too
type Event struct {
	Command Command
	Key string
	Value string
	Revision int64
}

type watcher struct {
	prefix string
	events chan Event
}

// Exposing channel to return events as read only
func (w *watcher) ReadEvents() <-chan Event {
	return w.events
}

func newWatcher(prefix string) *watcher {
	return &watcher{
		prefix: prefix,
		events: make(chan Event, 64),
	}
}

func (s *Store) Watch(prefix string) *watcher{
	nw := newWatcher(prefix)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watchers = append(s.watchers, nw)
	return nw
}
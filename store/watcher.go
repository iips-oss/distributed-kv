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
	Prefix string
	events chan Event
}

func NewWatcher(prefix string) *watcher {
	return &watcher{
		Prefix: prefix,
		events: make(chan Event, 64),
	}
}

func (s *Store) Watch(prefix string) *watcher{
	nw := NewWatcher(prefix)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Watchers = append(s.Watchers, nw)
	return nw
}
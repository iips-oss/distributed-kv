package store

// Using iota to define commands, as strings typos would compile silently
type Command int

const (
	SET Command = iota
	// GET
	DELETE
)

// Stringer interface for iota
// here case c==SET --- return "SET"
func (c Command) String() string {
	switch c{
	case SET:
		return "SET"
	case DELETE:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

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
	store *Store
}

// Exposing channel to return events as read only
func (w *watcher) ReadEvents() <-chan Event {
	return w.events
}

func newWatcher(prefix string, s *Store) *watcher {
	return &watcher{
		prefix: prefix,
		events: make(chan Event, 64),
		store: s,
	}
}

func (s *Store) Watch(prefix string, st *Store) *watcher{
	nw := newWatcher(prefix, st)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watchers = append(s.watchers, nw)
	return nw
}
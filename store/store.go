package store

type Store struct {
    data map[string]string
	wal *WAL
}

func NewStore(walPath string) (*Store, error) {
	w, err := OpenWAL(walPath)
	if err != nil {
		return nil, err
	}
	s := &Store{
		data: make(map[string]string),
		wal:  w,
	}
	if err := w.Replay(s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Put(key, value string) error {
	if err:=s.wal.Append(Put,key,value);err!=nil{
		return err
	}
    s.data[key] = value
	return nil
}

func (s *Store) Get(key string) (string, bool) {
    val, ok := s.data[key]
    return val, ok
}

func (s *Store) Delete(key string) error {
	if err:=s.wal.Append(Del,key,"");err!=nil{
		return err
	}
    delete(s.data, key)
	return nil
}
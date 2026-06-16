package store

import (
	"bufio"
	"encoding/json"
	"os"
)

type Entry struct {
	Op    string `json:"op"` // "PUT" or "DELETE"
	Key   string `json:"key"`
	Value string `json:"value"`
}
type WAL struct {
	file *os.File
}

// create a new WAL instance opening the file at the given path (creating it if it doesn't exist)
func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file: file,
	}, nil
}

// write a new entry to the WAL appending it to the file
func (w *WAL) Write(entry Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = w.file.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	return w.file.Sync()
}

// replay the WAL by reading all entries from the file and applying them using the provided function
func (w *WAL) Replay(apply func(Entry)) error {
	// go to beginning of file
	_, err := w.file.Seek(0, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(w.file)

	for scanner.Scan() {
		var entry Entry

		err := json.Unmarshal(scanner.Bytes(), &entry)
		if err != nil {
			continue
		}

		apply(entry)
	}

	return scanner.Err()
}

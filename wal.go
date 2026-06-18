package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const walFile = "wal.log"

// opens the file in append mode
// if file doesn't exist, creates it
func openWAL() (*os.File, error) {
	return os.OpenFile(walFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// os.O_APPEND → always write at the end, never overwrite
	// os.O_CREATE → create the file if it doesn't exist
	// os.O_WRONLY → open for writing only
	// 0644 → file permissions (owner can read/write, others can only read)
}

// writes a single command to the log file
func writeWAL(file *os.File, command string) error {
	_, err := file.WriteString(command + "\n")
	return err
}

// reads the log file line by line and replays into the store
func replayWAL(store *Store) error {
	file, err := os.Open(walFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "SET":
			store.mu.Lock()
			store.data[parts[1]] = parts[2]
			store.mu.Unlock()
		case "DEL":
			store.mu.Lock()
			delete(store.data, parts[1])
			store.mu.Unlock()
		case "FLUSH":
			store.mu.Lock()
			store.data = make(map[string]string)
			store.mu.Unlock()
		}
	}
	return nil
}
func readWAL() (string, error) {
	file, err := os.Open(walFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "WAL file is empty", nil
		}
		return "", err
	}
	defer file.Close()

	var result string
	lineNumber := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result += fmt.Sprintf("%d) %s\n", lineNumber, scanner.Text())
		lineNumber++
	}

	if result == "" {
		return "WAL file is empty", nil
	}
	return result, nil
}

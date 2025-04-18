package store

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type AOF struct {
	file *os.File
	mu   sync.Mutex
}

func NewAof(path string) *AOF {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	return &AOF{file: file}
}

func (a *AOF) AppendCmd(cmd string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.file.WriteString(fmt.Sprintf("%s\n", cmd))
}

func (a *AOF) Load() ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	file, err := os.Open(a.file.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

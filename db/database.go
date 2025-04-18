package db

import (
	"sync"
)

type Database struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]string),
	}
}

func (d *Database) Get(key string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	value, ok := d.data[key]
	return value, ok
}

func (d *Database) Set(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = value
}

func (d *Database) Del(key string) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.data[key]; ok {
		delete(d.data, key)
		return 1
	}
	return 0
}

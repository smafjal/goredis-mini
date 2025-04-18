package store

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
	data map[string]string
	ttl  map[string]time.Time
	aof  *AOF
	mu   sync.RWMutex
}

func NewDatabase(a *AOF) *Database {
	return &Database{
		data: make(map[string]string),
		ttl:  make(map[string]time.Time),
		aof:  a,
	}
}

func (d *Database) Get(key string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if exp, ok := d.ttl[key]; ok && time.Now().After(exp) {
		d.mu.RUnlock()
		d.mu.Lock()

		delete(d.data, key)
		delete(d.ttl, key)

		d.mu.Unlock()
		d.mu.RLock()
		return "", false
	}
	value, ok := d.data[key]
	return value, ok
}

func (d *Database) Set(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = value
	delete(d.ttl, key)
	if d.aof != nil {
		d.aof.AppendCmd(fmt.Sprintf("SET %s %s", key, value))
	}
}

func (d *Database) SetWithTTL(key, value string, ttl int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = value
	d.ttl[key] = time.Now().Add(time.Duration(ttl) * time.Second)
	if d.aof != nil {
		d.aof.AppendCmd(fmt.Sprintf("SETEX %s %d %s", key, ttl, value))
	}
}

func (d *Database) Del(key string) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.data[key]; ok {
		delete(d.data, key)
		delete(d.ttl, key)
		if d.aof != nil {
			d.aof.AppendCmd(fmt.Sprintf("DEL %s", key))
		}
		return 1
	}
	return 0
}

package pubsub

import (
	"fmt"
	"net"
	"sync"
)

type Pubsub struct {
	mu          sync.Mutex
	subscribers map[string][]net.Conn
}

func NewPubsub() *Pubsub {
	return &Pubsub{
		subscribers: make(map[string][]net.Conn),
	}
}

func (ps *Pubsub) Subscribe(channels []string, conn net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	fmt.Println("chs: ", channels)

	for _, ch := range channels {
		if _, ok := ps.subscribers[ch]; !ok {
			ps.subscribers[ch] = []net.Conn{}
		}
		ps.subscribers[ch] = append(ps.subscribers[ch], conn)
	}
}

func (ps *Pubsub) Publish(channel, msg string) int {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs, ok := ps.subscribers[channel]
	if !ok {
		return 0
	}
	for _, conn := range subs {
		pubMsg := fmt.Sprintf("$message: %s - %s\r\n", channel, msg)
		conn.Write([]byte(pubMsg))
	}
	return len(subs)
}

func (ps *Pubsub) Unsubscribe(channels []string, conn net.Conn) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, ch := range channels {
		conns, ok := ps.subscribers[ch]
		if !ok {
			continue
		}

		newConns := make([]net.Conn, 0, len(conns))
		for _, c := range conns {
			if c != conn {
				newConns = append(newConns, c)
			}
		}
		if len(newConns) == 0 {
			delete(ps.subscribers, ch)
		} else {
			ps.subscribers[ch] = newConns
		}
	}
}

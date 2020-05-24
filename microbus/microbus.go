package microbus

import "sync"

type MicroBus struct {
	mu     sync.RWMutex
	subs   map[string][]chan interface{}
	closed bool
}

func New() *MicroBus {
	mb := &MicroBus{}
	mb.subs = make(map[string][]chan interface{})
	return mb
}

func (mb *MicroBus) Subscribe(topic string) <-chan interface{} {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if mb.closed {
		return nil
	}

	ch := make(chan interface{}, 1)
	mb.subs[topic] = append(mb.subs[topic], ch)
	return ch
}

func (mb *MicroBus) Publish(topic string, content interface{}) {
	mb.mu.RLock()
	defer mb.mu.RUnlock()

	if mb.closed {
		return
	}

	for _, ch := range mb.subs[topic] {
		go func(ch chan interface{}) {
			ch <- content
		}(ch)
	}
}

func (mb *MicroBus) CloseAll() {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	if !mb.closed {
		mb.closed = true
		for _, subs := range mb.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

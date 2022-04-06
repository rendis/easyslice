package easyslice

import "sync"

type Closer struct {
	done   chan struct{}
	l      sync.RWMutex
	closed bool
}

func NewCloser() *Closer {
	return &Closer{
		done: make(chan struct{}),
	}
}

func (c *Closer) Done() <-chan struct{} {
	return c.done
}

func (c *Closer) Status() bool {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.closed
}

func (c *Closer) Close() {
	c.l.Lock()
	defer c.l.Unlock()
	if !c.closed {
		c.closed = true
		close(c.done)
	}
}

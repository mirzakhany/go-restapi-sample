package projectx

import (
	"context"
	"sync"
)

// Ctx contain project context
type Ctx struct {
	parent context.Context
	// This mutex protect Keys map
	mu sync.RWMutex

	// Keys is a key/value pair
	Keys map[string]interface{}
}

// New return a new instance of project context
func New(parent context.Context) *Ctx {
	return &Ctx{
		parent: parent,
		Keys:   make(map[string]interface{}),
	}
}

// Set is used to store a new key/value pair
func (c *Ctx) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Ctx) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// Done always returns nil (chan which will wait forever),
func (c *Ctx) Done() <-chan struct{} {
	return nil
}

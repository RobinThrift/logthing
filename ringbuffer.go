package main

import (
	"container/ring"
	"sync"
)

// A RingBuffer is a circular, thread safe (using a mutex) buffer.
type RingBuffer struct {
	ring  *ring.Ring
	first *ring.Ring
	m     *sync.Mutex
}

// NewRingBuffer creates a new empty ring buffer of size `size`
func NewRingBuffer(size int) *RingBuffer {
	ring := ring.New(size)
	return &RingBuffer{
		ring:  ring,
		first: ring,
		m:     &sync.Mutex{},
	}
}

// Insert a new value into the ring buffer. Thread safe thanks to mutex.
func (rb *RingBuffer) Insert(value interface{}) {
	rb.m.Lock()
	defer rb.m.Unlock()
	rb.ring.Value = value
	rb.ring = rb.ring.Next()
}

// Do loops over all the values in the ring buffer in order!
// The behavior of Do is undefined if f changes the ring buffer itself. Pls don't do that.
func (rb *RingBuffer) Do(f func(interface{})) {
	rb.m.Lock()
	defer rb.m.Unlock()
	rb.first.Do(f)
}

func (rb *RingBuffer) Clear() {
	for i := 0; i < rb.ring.Len(); i++ {
		rb.Insert(nil)
	}
}

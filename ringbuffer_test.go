package main

import (
	"testing"
)

func TestNewRingBuffer(t *testing.T) {
	t.Parallel()

	buffer := NewRingBuffer(10)
	if buffer == nil {
		t.Fatalf("Buffer is nil Oo")
	}
}

func TestRingBufferInsert(t *testing.T) {
	t.Parallel()

	buffer := NewRingBuffer(4)
	values := []rune{'a', 'b', 'c', 'd', 'e'}
	for _, v := range values {
		buffer.Insert(v)
	}

	values[0] = values[4] // we loop around once, and I'm lazy
	i := 0
	buffer.Do(func(v interface{}) {
		value, ok := v.(rune)
		if !ok {
			t.Fatalf("Value is not `rune`. Got %T", v)
		}

		if values[i] != value {
			t.Fatalf("Value should be %d. Got %d", values[i], value)
		}

		i++
	})
}

func TestThreadSaftey(t *testing.T) {
	t.Parallel()

	worker := func(value rune, buffer *RingBuffer, done chan rune) {
		buffer.Insert(value)
		done <- value
	}

	done := make(chan rune, 5)
	buffer := NewRingBuffer(4)
	values := []rune{'a', 'b', 'c', 'd', 'e'}
	for _, v := range values {
		go worker(v, buffer, done)
	}

	count := 0
	insertedValues := []rune{}

loop:
	for {
		select {
		case v := <-done:
			insertedValues = append(insertedValues, v)
			count++
			if count == 5 {
				break loop
			}
		}
	}

	insertedValues[0] = insertedValues[4] // we loop around once, and I'm lazy
	i := 0
	buffer.Do(func(v interface{}) {
		value, ok := v.(rune)
		if !ok {
			t.Fatalf("Value is not `rune`. Got %T", v)
		}

		if insertedValues[i] != value {
			t.Fatalf("Value should be %d. Got %d", insertedValues[i], value)
		}

		i++
	})
}

func TestRingBufferClear(t *testing.T) {
	t.Parallel()

	buffer := NewRingBuffer(4)
	values := []rune{'a', 'b', 'c', 'd', 'e'}
	for _, v := range values {
		buffer.Insert(v)
	}

	buffer.Clear()

	buffer.Do(func(v interface{}) {
		if v != nil {
			t.Fatalf("Value is not `nil`. Got %d", v)
		}
	})
}

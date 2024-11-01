package main

import (
	"fmt"
	"sync"
)

type NonBlockingQueue struct {
	queue chan interface{}
	lock  sync.Mutex
	cap   int
}

func NewNonBlockingQueue(cap int) *NonBlockingQueue {
	return &NonBlockingQueue{
		queue: make(chan interface{}, cap),
		cap:   cap,
	}
}

func (r *NonBlockingQueue) Length() int {
	return len(r.queue)
}

func (r *NonBlockingQueue) Capacity() int {
	return r.cap
}

func (r *NonBlockingQueue) Publish(val interface{}) bool {
	select {
	case r.queue <- val:
		return true
	default:
		return false
	}
}

func (r *NonBlockingQueue) Subscribe() (interface{}, bool) {
	select {
	case val := <-r.queue:
		return val, true
	default:
		return nil, false
	}
}

func main() {
	const cap = 10
	r := NewNonBlockingQueue(cap)

	for i := 0; i < cap; i++ {
		if r.Publish(i) {
			fmt.Println("Publish ", i, " to queue")
		} else {
			fmt.Println("Queue is full")
		}
	}

	for {
		if val, ok := r.Subscribe(); ok {
			fmt.Println("Subscribe value: ", val)
		} else {
			fmt.Println("Empty queue")
			break
		}
	}
}

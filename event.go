package gocache

import (
	"sync"
	"time"
)

type (
	event struct {
		key     interface{}
		timeout int64
	}

	eventlist struct {
		*ilist
		expireKey chan interface{}
		run       chan interface{}
		mu        sync.RWMutex
		len       int
	}

	ilist struct {
		head *node
	}

	node struct {
		e    event
		next *node
	}
)

func makeEventList() *eventlist {
	l := &eventlist{
		ilist: &ilist{
			head: &node{},
		},
		expireKey: make(chan interface{}),
		run:       make(chan interface{}),
	}
	go l.cleanExpire()
	return l
}

func (l *eventlist) orderInsert(e event) {
	l.mu.Lock()
	prev := l.head
	now := l.head.next
	for now != nil {
		if now.e.timeout > e.timeout {
			break
		}
		prev = now
		now = now.next
	}
	prev.next = &node{
		e:    e,
		next: now,
	}
	l.mu.Unlock()
	if l.len == 0 {
		l.run <- nil
	}
	l.len++
}

func (l *eventlist) cleanExpire() {
	for {
		var timeout time.Duration
		if l.len == 0 {
			<-l.run
		}
		l.mu.Lock()
		front := l.head.next
		timeout = time.Duration(front.e.timeout - time.Now().UnixNano())
		if timeout > 0 {
			l.mu.Unlock()
			time.Sleep(timeout - 2*time.Millisecond)
			l.mu.Lock()
		}
		l.expireKey <- front.e.key
		l.head.next = front.next
		l.len--
		l.mu.Unlock()
	}
}

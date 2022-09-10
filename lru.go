package gocache

import (
	"container/list"
	"sync"
)

type lrulist struct {
	List *list.List
	mu   sync.Mutex
}

func makeLrulist() *lrulist {
	return &lrulist{
		List: list.New(),
	}
}

func (l *lrulist) pushFront(entry entry) (e *list.Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e = l.List.PushFront(entry)
	return
}

func (l *lrulist) moveToFront(e *list.Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.List.MoveToFront(e)
}

func (l *lrulist) remove(e *list.Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.List.Remove(e)
}

func (l *lrulist) back() (e *list.Element) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e = l.List.Back()
	return
}

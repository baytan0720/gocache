package gocache

import (
	"container/list"
	"sync"
)

type lrulist struct {
	*list.List
	mu sync.Mutex
}

func makeLrulist() *lrulist {
	return &lrulist{
		List: list.New(),
	}
}

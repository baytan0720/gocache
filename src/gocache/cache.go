package gocache

import "sync"

type cache struct {
	memData
	lrulist
	lrumap map[interface{}]*node
	cap    int
	mu     sync.Mutex
	closed bool
}

func New(cap ...int) *cache {
	c := &cache{
		memData: *makeMemData(),
		lrulist: *makeLrulist(),
		lrumap:  make(map[interface{}]*node),
	}
	if len(cap) > 0 {
		c.cap = cap[0]
	}
	return c
}

func (c *cache) SetCap(cap int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cap = cap
	if cap < 1 {
		return
	}
	for c.Size() > cap {
		key := c.deletelast()
		c.remove(key)
		delete(c.lrumap, key)
	}
}

func (c *cache) Set(key, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, val)
	node := &node{
		data: key,
	}
	c.insertHead(node)
	c.lrumap[key] = node
	if c.cap > 0 && c.Size() > c.cap {
		key := c.deletelast()
		c.remove(key)
		delete(c.lrumap, key)
	}
}

func (c *cache) SetIfNotExist(key, val interface{}) {
	_, ok := c.get(key)
	if ok {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(key, val)
	node := &node{
		data: key,
	}
	c.insertHead(node)
	c.lrumap[key] = node
	if c.cap > 0 && c.Size() > c.cap {
		key := c.deletelast()
		c.remove(key)
		delete(c.lrumap, key)
	}
}

func (c *cache) Get(key interface{}) (val interface{}, ok bool) {
	val, ok = c.get(key)
	if !ok {
		return val, ok
	}
	c.mu.Lock()
	node := c.lrumap[key]
	c.mu.Unlock()
	c.delete(node)
	c.insertHead(node)
	return
}

func (c *cache) GetorSet(key, val interface{}) (v interface{}, ok bool) {
	v, ok = c.get(key)
	if !ok {
		c.Set(key, val)
		return val, ok
	}
	c.mu.Lock()
	node := c.lrumap[key]
	c.mu.Unlock()
	c.delete(node)
	c.insertHead(node)
	return
}

func (c *cache) Del(key ...interface{}) {
	keys := c.remove(key...)
	c.mu.Lock()
	for _, v := range keys {
		c.delete(c.lrumap[v])
		delete(c.lrumap, v)
	}
	c.mu.Unlock()
}

func (c *cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clear()
	c.memData = *makeMemData()
	c.lrulist = *makeLrulist()
	c.lrumap = make(map[interface{}]*node)
}

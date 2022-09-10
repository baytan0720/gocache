package gocache

import "time"

type cache struct {
	*data
	lrulist   *lrulist
	eventlist *eventlist
	cap       int
}

type any = interface{}

// New接收一个cap，不传入或者小于1则为无限制容量，返回一个cache
func New(cap ...int) *cache {
	c := &cache{
		data:      makeData(),
		lrulist:   makeLrulist(),
		eventlist: makeEventList(),
	}
	if len(cap) > 0 {
		c.cap = cap[0]
	}
	go func(expire chan interface{}) {
		for {
			key := <-expire
			c.Del(key)
		}
	}(c.eventlist.expireKey)
	return c
}

// SetCap修改cap，为了安全，新cap不得小于c.Size(),小于1除外
func (c *cache) SetCap(cap int) {
	if cap < 1 {
		c.cap = cap
		return
	}
	if cap < c.Size() {
		return
	} else {
		c.cap = cap
	}
}

// Set接收任意类型的key和val，并将其写入cache
func (c *cache) Set(key, val any) {
	Entry := entry{
		key:     key,
		val:     val,
		timeout: -1,
	}
	c.data.add(key, c.lrulist.pushFront(Entry))
	if c.cap > 0 && c.Size() > c.cap {
		key := c.lrulist.back().Value.(entry).key
		c.Del(key)
		return
	}
}

// SetWithTimeout接收任意类型的key和val以及超时时间，并将其写入cache，到期后删除
func (c *cache) SetWithTimeout(key, val any, timeout time.Duration) {
	expire := time.Now().UnixNano() + timeout.Nanoseconds()
	Entry := entry{
		key:     key,
		val:     val,
		timeout: expire,
	}
	c.data.add(key, c.lrulist.pushFront(Entry))
	c.eventlist.orderInsert(event{
		key:     key,
		timeout: expire,
	})
	if c.cap > 0 && c.Size() > c.cap {
		key := c.lrulist.back().Value.(entry).key
		c.Del(key)
		return
	}
}

// SetIfNotExist接收任意类型的key和val，当key存在时返回false；当key不存在时将其写入cache，并返回true
func (c *cache) SetIfNotExist(key, val any) (ok bool) {
	if _, ok := c.data.get(key); ok {
		return false
	}
	c.Set(key, val)
	return true
}

// Get接收任意类型的key，如果key存在则返回val和true，否则返回nil和false
func (c *cache) Get(key any) (val interface{}, ok bool) {
	e, ok := c.data.get(key)
	if !ok {
		return nil, ok
	}
	entry := e.Value.(entry)
	c.lrulist.moveToFront(e)
	return entry.val, ok
}

// GetTimeout接收任意类型的key，如果key存在则返回剩余过期时间，如果key是无超时则返回-1；如果不存在key则返回false
func (c *cache) GetTimeOut(key any) (timeout time.Duration, ok bool) {
	e, ok := c.data.get(key)
	if !ok {
		return -1, false
	}
	entry := e.Value.(entry)
	return time.Duration(entry.timeout - time.Now().UnixNano()), ok
}

// GetOrSet接收任意类型的key和val，如果key存在则返回val，不存在则写入cache，并返回val
func (c *cache) GetOrSet(key, val any) (Val interface{}) {
	if Val, ok := c.Get(key); ok {
		return Val
	}
	c.Set(key, val)
	return Val
}

// Contains接收任意类型的key，如果key存在则返回true，否则返回false
func (c *cache) Contains(key any) bool {
	_, ok := c.data.get(key)
	return ok
}

// Del接收一个或多个key，并将其删除
func (c *cache) Del(key ...any) {
	es := c.data.remove(key...)
	for _, e := range es {
		c.lrulist.remove(e)
	}
}

// Cap返回cache的容量
func (c *cache) Cap() int {
	return c.cap
}

// Clear用以清空cache，谨慎使用
func (c *cache) Clear() {
	c.data = makeData()
	c.lrulist = makeLrulist()
	c.eventlist = makeEventList()
}

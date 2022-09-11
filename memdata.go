package gocache

import (
	"container/list"
	"sync"
)

type entry struct {
	key     interface{}
	val     interface{}
	timeout int64
}

type data struct {
	mu   sync.RWMutex
	data map[interface{}]*list.Element
}

func makeData() *data {
	return &data{
		data: make(map[interface{}]*list.Element),
	}
}

func (d *data) add(key interface{}, entry *list.Element) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = entry
}

func (d *data) get(key interface{}) (e *list.Element, ok bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	e, ok = d.data[key]
	if !ok {
		return nil, ok
	}
	return
}

func (d *data) remove(key ...interface{}) (es []*list.Element) {
	d.mu.Lock()
	defer d.mu.Unlock()
	es = make([]*list.Element, 0, len(key))
	for _, key := range key {
		if e, ok := d.data[key]; ok {
			es = append(es, e)
			delete(d.data, key)
		}
	}
	return
}

// Size返回cache中已缓存的数量
func (d *data) Size() int {
	return len(d.data)
}

// Entrys以键值对的形式返回cache中所有缓存对象
func (d *data) Entrys() (entrys [][]interface{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	entrys = make([][]interface{}, 0, d.Size())
	for _, v := range d.data {
		e := v.Value.(entry)
		entrys = append(entrys, []interface{}{e.key, e.val})
	}
	return
}

// Keys返回cache中所有key
func (d *data) Keys() (keys []interface{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	keys = make([]interface{}, 0, d.Size())
	for k := range d.data {
		keys = append(keys, k)
	}
	return
}

// Vals返回cache中所有val
func (d *data) Vals() (vals []interface{}) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	vals = make([]interface{}, 0, d.Size())
	for _, v := range d.data {
		vals = append(vals, v.Value.(entry).val)
	}
	return
}

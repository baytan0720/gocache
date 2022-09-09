package gocache

import (
	"sync"
	"time"
)

// type memKey struct {
// 	key     interface{}
// 	timeout time.Duration
// }

type memVal struct {
	val     interface{}
	timeout time.Duration
}

type memData struct {
	mu   sync.RWMutex
	data map[interface{}]memVal
}

func makeMemData() *memData {
	return &memData{
		data: make(map[interface{}]memVal),
	}
}

func (d *memData) Update(key, val interface{}) (oldVal interface{}, exist bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if v, ok := d.data[key]; ok {
		d.data[key] = memVal{
			val: val,
		}
		return v.val, true
	}
	return nil, false
}

func (d *memData) UpdateTimeout(key interface{}, timeout time.Duration) (oldtimeout time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if v, ok := d.data[key]; ok {
		d.data[key] = memVal{
			timeout: timeout,
		}
		return v.timeout
	}
	return -1
}

func (d *memData) remove(key ...interface{}) (removedKeys []interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	removedKeys = make([]interface{}, 0, len(key))
	for _, k := range key {
		if _, ok := d.data[k]; ok {
			removedKeys = append(removedKeys, k)
			delete(d.data, k)
		}
	}
	return removedKeys
}

func (d *memData) Datas() map[interface{}]interface{} {
	d.mu.RLock()
	m := make(map[interface{}]interface{}, len(d.data))
	for k, v := range d.data {
		m[k] = v.val
	}
	d.mu.RUnlock()
	return m
}

func (d *memData) Keys() []interface{} {
	d.mu.RLock()
	var (
		index = 0
		keys  = make([]interface{}, len(d.data))
	)
	for k := range d.data {
		keys[index] = k
		index++
	}
	d.mu.RUnlock()
	return keys
}

func (d *memData) Values() []interface{} {
	d.mu.RLock()
	var (
		index  = 0
		values = make([]interface{}, len(d.data))
	)
	for _, v := range d.data {
		values[index] = v.val
		index++
	}
	d.mu.RUnlock()
	return values
}

func (d *memData) Size() (size int) {
	d.mu.RLock()
	size = len(d.data)
	d.mu.RUnlock()
	return size
}

func (d *memData) clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data = make(map[interface{}]memVal)
}

func (d *memData) get(key interface{}) (val interface{}, ok bool) {
	d.mu.RLock()
	v, ok := d.data[key]
	d.mu.RUnlock()
	if !ok {
		return nil, ok
	}
	return v.val, ok
}

func (d *memData) set(key interface{}, val interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.data[key] = memVal{
		val: val,
	}
}

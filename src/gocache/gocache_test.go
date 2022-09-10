package gocache

import "testing"

func TestList(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		go c.Set(i, i)
	}
	c.SetCap(5)
}

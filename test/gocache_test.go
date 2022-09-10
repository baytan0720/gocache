package test

import (
	"gocache"
	"testing"
)

func TestGocache(t *testing.T) {
	c := gocache.New()
	for i := 0; i < 10000000; i++ {
		go c.Set(i, i)
	}
}

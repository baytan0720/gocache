package gocache_test

import (
	"gocache/src/gocache"
	"testing"
)

func TestSet(t *testing.T) {
	c := gocache.New(-1)
	c.Set(1, 1)
}

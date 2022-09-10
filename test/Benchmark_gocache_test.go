package test

import (
	"gocache"
	"testing"
)

func Benchmark_Set(b *testing.B) {
	c := gocache.New()
	i := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(i, i)
			i++
		}
	})
}

func Benchmark_Get(b *testing.B) {
	c := gocache.New()
	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
	i := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get(i)
			i++
		}
	})
}

func Benchmark_Del(b *testing.B) {
	c := gocache.New()
	for i := 0; i < b.N; i++ {
		c.Set(i, i)
	}
	i := 0
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Del(i)
			i++
		}
	})
}

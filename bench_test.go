package main

// $ go test -bench=. -cpuprofile=profile && go tool pprof profile

import (
	"testing"
)

func BenchmarkCache(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Cache()
	}
}

// 测试并发效率
func BenchmarkLoopsParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) { //并发
		for pb.Next() {
			Cache()
		}
	})
}

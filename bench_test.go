package main

// $ go test -bench=. -cpuprofile=profile && go tool pprof profile

import "testing"

func BenchmarkFib20(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Cache()
	}
}

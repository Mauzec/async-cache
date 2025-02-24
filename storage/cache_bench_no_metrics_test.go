package storage

import "testing"

const parallelFactor = 10_000

func Benchmark_RWMutex_BalanceLoad(b *testing.B) {
	// b.Skip()
	c := NewAsyncCache()
	for i := 0; i < b.N; i++ {
		emulateLoadBench(c, parallelFactor)
	}
}

func Benchmark_Mutex_BalanceLoad(b *testing.B) {
	// b.Skip()
	c := NewToughtAsyncCache()
	for i := 0; i < b.N; i++ {
		emulateLoadBench(c, parallelFactor)
	}
}

package benchmark

import "testing"

// Run basic benchmarks by
// $ go test --bench .

func BenchmarkHeavyJob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = HeavyJob(i)
	}
}

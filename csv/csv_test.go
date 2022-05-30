package csv

import "testing"

func BenchmarkProcessFile(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ProcessFile()
	}
}

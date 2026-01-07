package ratelimit

import (
	"testing"
	"time"
)

// --------------------------------------------------------------------

func BenchmarkLimit(b *testing.B) {
	rl := New(1000, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Limit()
	}
}

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	// 	RegisterFailHandler(Fail)
	// RunSpecs(t, "github.com/bsm/ratelimit")

	t.Log("\nshould accurately rate-limit at small rates")
	var count int
	rl := New(10, time.Minute)
	for !rl.Limit() {
		count++
	}
	if count != 10 {
		t.Fail()
	}
}

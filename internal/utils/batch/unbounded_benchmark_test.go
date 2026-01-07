package batch

import "testing"

func BenchmarkChannel(b *testing.B) {
	done := make(chan struct{})
	ch := make(chan int)
	go func() {
		defer close(done)
		for range ch {
		}
	}()
	for i := 0; i < b.N; i++ {
		ch <- i
	}
	close(ch)
	<-done
}

func BenchmarkUnboundedChannel(b *testing.B) {
	done := make(chan struct{})
	tx, rx := NewUnboundedChannel[int](nil)
	go func() {
		defer close(done)
		for range rx {
		}
	}()
	for i := 0; i < b.N; i++ {
		tx <- i
	}
	close(tx)
	<-done
}

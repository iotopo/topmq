package batch

import (
	"testing"
	"time"
)

// tx 发送至到 rx 接收到值有一定延迟
// tx 关闭到 rx 关闭有一定延迟
var testUnboundedChannelDelay = 100 * time.Millisecond

func assertUnboundedChannelEmpty(t *testing.T, rx <-chan int) {
	select {
	case val, ok := <-rx:
		if ok {
			t.Fatalf("assert channel empty, got value %d", val)
		}
		t.Fatal("channel is closed")
	case <-time.After(testUnboundedChannelDelay):
		return
	}
}

func assertUnboundedChannelEmptyAndClosed(t *testing.T, rx <-chan int) {
	select {
	case val, ok := <-rx:
		if ok {
			t.Fatalf("assert channel empty, got value %d", val)
		}
	case <-time.After(testUnboundedChannelDelay):
		t.Fatal("assert channel closed")
	}
}

func assertUnboundedChannelRx(t *testing.T, rx <-chan int, expected int) {
	select {
	case val, ok := <-rx:
		if !ok {
			t.Fatal("channel is closed")
		}
		if val != expected {
			t.Fatalf("expect %d, got %d", expected, val)
		}
	case <-time.After(testUnboundedChannelDelay):
		t.Fatal("channel is empty")
	}
}

func TestUnboundedChannelEmpty(t *testing.T) {
	tx, rx := NewUnboundedChannel[int](nil)
	assertUnboundedChannelEmpty(t, rx)
	close(tx)
	time.Sleep(testUnboundedChannelDelay)
	assertUnboundedChannelEmptyAndClosed(t, rx)
}

func TestUnboundedChannelPushAndPop(t *testing.T) {
	tx, rx := NewUnboundedChannel[int](nil)

	tx <- 1001
	tx <- 1002
	assertUnboundedChannelRx(t, rx, 1001)
	tx <- 1003
	close(tx)
	assertUnboundedChannelRx(t, rx, 1002)
	assertUnboundedChannelRx(t, rx, 1003)
	assertUnboundedChannelEmptyAndClosed(t, rx)

}

func TestUnboundedChannelWithStopping(t *testing.T) {
	stopping := make(chan struct{})
	tx, rx := NewUnboundedChannel[int](stopping)

	tx <- 1001
	tx <- 1002
	assertUnboundedChannelRx(t, rx, 1001)
	tx <- 1003
	close(stopping)
	time.Sleep(testUnboundedChannelDelay)
	assertUnboundedChannelEmptyAndClosed(t, rx)
	close(tx)
}

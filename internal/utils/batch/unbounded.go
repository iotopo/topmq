package batch

type unboundedChannel[T any] struct {
	queue *Queue[T]
	in    <-chan T
	out   chan<- T
	kill  chan struct{}
}

// 由发送端和接受端组成的一个无限缓冲的 channel
// 使用结束后应关闭发送端和 kill，否则会造成 goroutine 泄露
// 如果 kill != nil，当 kill 关闭时，接收端也会尽快关闭，未接受的数据被丢弃
// 如果 kill == nil，当发送端关闭时，接收端会等待数据全部接收完毕后再关闭
//
//goland:noinspection GoUnusedExportedFunction
func NewUnboundedChannel[T any](kill chan struct{}) (chan<- T, <-chan T) {
	in := make(chan T)
	out := make(chan T)
	ch := unboundedChannel[T]{
		queue: NewQueue[T](),
		in:    in,
		out:   out,
		kill:  kill,
	}
	go ch.waitIn()
	return in, out
}

func (ch unboundedChannel[T]) waitIn() {
	for {
		select {
		case <-ch.kill:
			close(ch.out)
			return
		case elem, ok := <-ch.in:
			if ok {
				if ch.waitBoth(elem) {
					return
				}
				// out 空闲，继续等待 in
			} else {
				// in 被关闭，因为 out 空闲，直接关闭即可
				close(ch.out)
				return
			}
		}
	}
}

// 等待 in 和 out，参数 elem 是当前要发给 out
// 返回 false 表示 elem 已发送，out 空闲（queue为空）
// 返回 true 表示 out 和 queue 已关闭（elem已发送或丢弃）
func (ch unboundedChannel[T]) waitBoth(elem T) bool {
	var next T
	var ok bool
	for {
		select {
		case <-ch.kill:
			close(ch.out)
			return true
		case ch.out <- elem:
			// out 的速度大于 in，从 queue 中补充
			elem, ok = ch.queue.PopFront()
			// 无元素可补充，out 空闲，继续等待 in 即可
			if !ok {
				return false
			}
		case next, ok = <-ch.in:
			if ok {
				// in 的速度大于 out，添加到队列，继续等待 in 和 out
				ch.queue.PushBack(next)
			} else {
				// in 已关闭，结束
				ch.waitOutAndClose(elem)
				return true
			}
		}
	}
}

// 发送 elem 和 queue 中剩余元素到 out，然后关闭 out
func (ch unboundedChannel[T]) waitOutAndClose(elem T) {
	defer close(ch.out)

	ok := true
	for ok {
		select {
		case <-ch.kill:
			return
		case ch.out <- elem:
			elem, ok = ch.queue.PopFront()
		}
	}
}

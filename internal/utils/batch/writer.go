package batch

import (
	"github.com/sirupsen/logrus"
	"time"
)

type WriterOptions[T any] struct {
	WriteFunc     func([]T)     // required
	BatchSize     int           // required
	BatchInterval time.Duration // default: 1s
	// 如果 AdditionalBufferSize >= 0，缓冲区大小为 BatchSize + AdditionalBufferSize
	//   当缓冲区满且写操作正在执行时，Append 将会阻塞
	// 如果 AdditionalBufferSize < 0，使用无限缓冲区，Append 不会阻塞
	//   但注意如果 WriteFunc 的速度如果始终跟不上，内存可能会无限增长
	AdditionalBufferSize int
}

type Writer[T any] struct {
	inputChan chan<- T
	kill      chan struct{}
	done      chan struct{}
}

func NewBatchWriter[T any](batchSize int, batchTimeout time.Duration, writeFunc func(datas []T)) *Writer[T] {
	return NewWriter(WriterOptions[T]{
		WriteFunc:     writeFunc,
		BatchSize:     batchSize,
		BatchInterval: batchTimeout,
	})
}

//goland:noinspection GoUnusedExportedFunction
func NewWriter[T any](opts WriterOptions[T]) *Writer[T] {
	if opts.BatchInterval < time.Millisecond*10 {
		logrus.Warnf("batch interval too small, reset to 1s")
		opts.BatchInterval = time.Second
	}

	kill := make(chan struct{})
	done := make(chan struct{})

	if opts.AdditionalBufferSize < 0 {
		tx, rx := NewUnboundedChannel[T](kill)
		go batchWrite(opts.BatchSize, opts.BatchInterval, rx, opts.WriteFunc, done, kill)
		return &Writer[T]{
			inputChan: tx,
			kill:      kill,
			done:      done,
		}
	} else {
		var inputChan chan T
		if opts.AdditionalBufferSize >= 0 {
			inputChan = make(chan T, opts.BatchSize+opts.AdditionalBufferSize)
		} else {
			inputChan = make(chan T)
		}
		go batchWrite(opts.BatchSize, opts.BatchInterval, inputChan, opts.WriteFunc, done, kill)
		return &Writer[T]{
			inputChan: inputChan,
			kill:      kill,
			done:      done,
		}
	}
}

// 如果 Close 被调用，Append 会立即返回 false 并丢弃数据
func (w *Writer[T]) Append(data T) bool {
	// NOTE: 如果 w.kill 和 w.inputChan 都已经被关闭
	// 这里不会返回 false 而是会 panic
	select {
	case <-w.kill:
		return false
	default:
		select {
		case w.inputChan <- data:
			return true
		case <-w.kill:
			return false
		}
	}
}

// NOTE: 在已关闭的 Writer 上调用 Close 会 panic
func (w *Writer[T]) Close() {
	close(w.kill)
	// NOTE: 我们不关闭 inputChan，因为会导致 Append() 在 Close() 之后调用时产生 panic
	// 为此，batchWrite 的 kill 参数不能为 nil，只使用 kill 来退出 goroutine
	// close(w.inputChan)
	<-w.done
}

func batchWrite[T any](
	batchSize int,
	batchInterval time.Duration,
	inputChan <-chan T,
	writeFunc func(buf []T),
	done chan struct{},
	kill chan struct{},
) {
	defer close(done)

	buf := make([]T, 0, batchSize)
	ticker := time.NewTicker(batchInterval)

	for {
		select {
		case <-kill:
			ticker.Stop()
			return
		case data, ok := <-inputChan:
			if ok {
				buf = append(buf, data)
				if len(buf) >= batchSize {
					writeFunc(buf)
					buf = buf[:0]
				}
			} else {
				ticker.Stop()
				// 如果 kill != nil，丢弃剩余数据
				if kill == nil && len(buf) > 0 {
					writeFunc(buf)
				}
				return
			}
		case <-ticker.C:
			if len(buf) > 0 {
				writeFunc(buf)
				buf = buf[:0]
			}
		}
	}
}

package batch

import "testing"

// 用 slice 实现的队列，当数据规模较小时，性能很好
// 有一个很大的缺陷就是内存增长后不会回收
type SliceQueue[T any] struct {
	slice []T
}

func (q *SliceQueue[T]) PushBack(elem T) {
	q.slice = append(q.slice, elem)
}

func (q *SliceQueue[T]) PopFront() (elem T, ok bool) {
	if len(q.slice) > 0 {
		elem = q.slice[0]
		ok = true
	}
	return
}

// 单链表实现的队列，性能很差，优点是实现简单
type singleLinkedQueueItem[T any] struct {
	value T
	next  *singleLinkedQueueItem[T]
}

type SingleLinkedQueue[T any] struct {
	head *singleLinkedQueueItem[T]
	tail *singleLinkedQueueItem[T]
}

func (q *SingleLinkedQueue[T]) PushBack(elem T) {
	b := &singleLinkedQueueItem[T]{value: elem}

	if q.tail == nil {
		q.tail = b
		q.head = b
	} else {
		q.tail.next = b
		q.tail = b
	}
}

func (q *SingleLinkedQueue[T]) PopFront() (elem T, ok bool) {
	if q.head == nil {
		return
	}
	elem = q.head.value
	ok = true
	q.head = q.head.next
	if q.head == nil {
		q.tail = nil
	}
	return
}

func NewSingleLinkedQueue[T any]() *SingleLinkedQueue[T] {
	return &SingleLinkedQueue[T]{}
}

type IQueue[T any] interface {
	PushBack(T)
	PopFront() (T, bool)
}

func benchmarkQueueImpl[Q IQueue[int]](b *testing.B, queue Q) {
	for i := 0; i < b.N; i++ {
		queue.PushBack(i)
	}
	for i := 0; i < b.N; i++ {
		queue.PopFront()
	}
}

func BenchmarkSliceQueue(b *testing.B) {
	benchmarkQueueImpl(b, &SliceQueue[int]{})
}

func BenchmarkSingleLinkedQueue(b *testing.B) {
	benchmarkQueueImpl(b, NewSingleLinkedQueue[int]())
}

func BenchmarkQueue(b *testing.B) {
	benchmarkQueueImpl(b, NewQueue[int]())
}

package batch

const queueBlockSize = 32

type queueBlock[T any] struct {
	next   *queueBlock[T]
	values [queueBlockSize]T
	index  uint8
	size   uint8
}

type Queue[T any] struct {
	head *queueBlock[T]
	tail *queueBlock[T]
}

// NewQueue 创建一个链表实现的队列
// 为了减少内存分配次数，内部没有直接将元素组成链表，而是先组成 queueBlock，再组成链表
//goland:noinspection GoUnusedExportedFunction
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) PushBack(elem T) {
	if q.tail == nil {
		b := new(queueBlock[T])
		b.values[0] = elem
		b.size = 1
		q.tail = b
		q.head = b
		return
	}

	if q.tail.size >= queueBlockSize {
		b := new(queueBlock[T])
		b.values[0] = elem
		b.size = 1
		q.tail.next = b
		q.tail = b
		return
	}

	q.tail.values[q.tail.size] = elem
	q.tail.size++
}

func (q *Queue[T]) Front() *T {
	if q.head == nil || q.head.index >= q.head.size {
		return nil
	}
	return &q.head.values[q.head.index]
}

// NOTE: 这个方法是 O(n) 的，如果需要 O(1) 的，可以再封装一层，在 PushBack 和 PopFront 时记录
func (q *Queue[T]) Len() (size int) {
	for b := q.head; b != nil; b = b.next {
		size += int(b.size - b.index)
	}
	return
}

func (q *Queue[T]) IsEmpty() bool {
	return q.head == nil || q.head.index >= q.head.size
}

func (q *Queue[T]) PopFront() (elem T, ok bool) {
	if q.head == nil || q.head.index >= q.head.size {
		return
	}
	elem = q.head.values[q.head.index]
	ok = true
	q.head.index++
	if q.head.index >= q.head.size {
		if q.head.size >= queueBlockSize {
			q.head = q.head.next
			if q.head == nil {
				q.tail = nil
			}
		} else {
			// 如果 block 未使用满，保留此 block
			// 这种情况只会发生在只有一个 block 时
			// 此时可以重置 index 和 size
			q.head.index = 0
			q.head.size = 0
		}
	}
	return
}

package batch

import "testing"

func assertQueueLen(t *testing.T, queue *Queue[int], expected int) {
	if queue.Len() != expected {
		t.Fatalf("assert Len() == %d", expected)
	}
	if expected > 0 {
		if queue.IsEmpty() {
			t.Fatal("assert not IsEmpty()")
		}
	} else {
		if !queue.IsEmpty() {
			t.Fatal("assert IsEmpty()")
		}
	}
}

func assertQueueEmpty(t *testing.T, queue *Queue[int]) {
	if queue.Front() != nil {
		t.Fatal("assert Front() == nil")
	}
	if _, ok := queue.PopFront(); ok {
		t.Fatal("assert PopFront() == (_, false)")
	}
}

func assertQueueFrontAndPop(t *testing.T, queue *Queue[int], expected int) {
	if front := queue.Front(); front == nil || *front != expected {
		t.Fatalf("assert Front() == %d", expected)
	}
	if front, ok := queue.PopFront(); !ok || front != expected {
		t.Fatalf("assert PopFront() == (%d, true)", expected)
	}
}

func TestQueueEmpty(t *testing.T) {
	queue := NewQueue[int]()
	assertQueueLen(t, queue, 0)
	assertQueueEmpty(t, queue)
}

func TestQueuePushAndPop(t *testing.T) {
	queue := NewQueue[int]()

	queue.PushBack(1001)
	assertQueueLen(t, queue, 1)
	assertQueueFrontAndPop(t, queue, 1001)
	assertQueueLen(t, queue, 0)
	assertQueueEmpty(t, queue)

	queue.PushBack(1002)
	assertQueueLen(t, queue, 1)
	assertQueueFrontAndPop(t, queue, 1002)
	assertQueueLen(t, queue, 0)
	assertQueueEmpty(t, queue)
}

func TestQueuePushMultiAndPop(t *testing.T) {
	queue := NewQueue[int]()

	queue.PushBack(1001)
	queue.PushBack(1002)
	assertQueueLen(t, queue, 2)
	assertQueueFrontAndPop(t, queue, 1001)
	queue.PushBack(1003)
	assertQueueLen(t, queue, 2)
	assertQueueFrontAndPop(t, queue, 1002)
	assertQueueLen(t, queue, 1)
	assertQueueFrontAndPop(t, queue, 1003)
	assertQueueLen(t, queue, 0)
	assertQueueEmpty(t, queue)
}

func TestQueuePushMultiBlockAndPop(t *testing.T) {
	queue := NewQueue[int]()

	for i := 1; i <= 70; i++ {
		queue.PushBack(1100 + i)
		assertQueueLen(t, queue, i)
	}
	for i := 1; i <= 70; i++ {
		assertQueueLen(t, queue, 71-i)
		assertQueueFrontAndPop(t, queue, 1100+i)
	}
	assertQueueLen(t, queue, 0)
	assertQueueEmpty(t, queue)
}

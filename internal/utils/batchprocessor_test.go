package utils

import (
	"sync"
	"testing"
	"time"
)

func TestBatchProcessor(t *testing.T) {
	// 测试基本功能
	t.Run("basic functionality", func(t *testing.T) {
		var processedItems []int
		var mu sync.Mutex
		bp := NewBatchProcessor[int](3, func(items []int) {
			mu.Lock()
			defer mu.Unlock()
			processedItems = append(processedItems, items...)
		})

		bp.Start()
		// 添加5个元素，应该触发两次批处理
		for i := 0; i < 5; i++ {
			bp.Add(i)
		}

		// 等待一小段时间确保处理完成
		time.Sleep(200 * time.Millisecond)
		bp.Stop()

		mu.Lock()
		defer mu.Unlock()
		if len(processedItems) != 5 {
			t.Errorf("expected 5 items, got %d", len(processedItems))
		}
		// 验证处理的顺序
		for i := 0; i < 5; i++ {
			if i < len(processedItems) && processedItems[i] != i {
				t.Errorf("item at index %d: expected %d, got %d", i, i, processedItems[i])
			}
		}
	})

	// 测试停止功能
	t.Run("stop functionality", func(t *testing.T) {
		var processedItems []int
		var mu sync.Mutex
		bp := NewBatchProcessor[int](3, func(items []int) {
			mu.Lock()
			defer mu.Unlock()
			processedItems = append(processedItems, items...)
		})

		bp.Start()

		// 添加一些元素
		for i := 0; i < 2; i++ {
			bp.Add(i)
		}

		// 等待一小段时间确保处理完成
		time.Sleep(200 * time.Millisecond)
		bp.Stop()

		// 尝试添加更多元素
		bp.Add(3)

		mu.Lock()
		defer mu.Unlock()
		if len(processedItems) != 2 {
			t.Errorf("expected 2 items, got %d", len(processedItems))
		}
	})

	// 测试并发添加
	t.Run("concurrent add", func(t *testing.T) {
		var processedItems []int
		var mu sync.Mutex
		bp := NewBatchProcessor[int](10, func(items []int) {
			mu.Lock()
			defer mu.Unlock()
			processedItems = append(processedItems, items...)
		})

		bp.Start()

		// 并发添加100个元素
		var wg sync.WaitGroup
		errChan := make(chan error, 100)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(start int) {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					bp.Add(start + j)
				}
			}(i * 10)
		}

		wg.Wait()
		close(errChan)

		// 检查是否有错误发生
		for err := range errChan {
			t.Errorf("error during concurrent add: %v", err)
		}

		// 等待一小段时间确保处理完成
		time.Sleep(200 * time.Millisecond)
		bp.Stop()

		mu.Lock()
		defer mu.Unlock()
		if len(processedItems) != 100 {
			t.Errorf("expected 100 items, got %d", len(processedItems))
		}
	})

	// 测试批处理大小
	t.Run("batch size", func(t *testing.T) {
		var batchSizes []int
		var mu sync.Mutex
		bp := NewBatchProcessor[int](3, func(items []int) {
			mu.Lock()
			defer mu.Unlock()
			batchSizes = append(batchSizes, len(items))
		})

		bp.Start()
		// 添加7个元素，应该触发3次批处理（3+3+1）
		for i := 0; i < 7; i++ {
			bp.Add(i)
		}

		// 等待一小段时间确保处理完成
		time.Sleep(200 * time.Millisecond)
		bp.Stop()

		mu.Lock()
		defer mu.Unlock()
		if len(batchSizes) != 3 {
			t.Errorf("expected 3 batches, got %d", len(batchSizes))
		}
		if len(batchSizes) == 3 && (batchSizes[0] != 3 || batchSizes[1] != 3 || batchSizes[2] != 1) {
			t.Errorf("unexpected batch sizes: %v", batchSizes)
		}
	})

	// 测试定时刷新
	t.Run("timer flush", func(t *testing.T) {
		var processedItems []int
		var mu sync.Mutex
		bp := NewBatchProcessor[int](3, func(items []int) {
			mu.Lock()
			defer mu.Unlock()
			processedItems = append(processedItems, items...)
		})

		bp.Start()
		// 添加2个元素，应该在定时器触发时被处理
		for i := 0; i < 2; i++ {
			bp.Add(i)
		}

		// 等待定时器触发
		time.Sleep(200 * time.Millisecond)
		bp.Stop()

		mu.Lock()
		defer mu.Unlock()
		if len(processedItems) != 2 {
			t.Errorf("expected 2 items to be processed by timer, got %d", len(processedItems))
		}
	})
}

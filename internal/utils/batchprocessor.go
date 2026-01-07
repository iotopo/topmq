package utils

import (
	"sync"
	"time"
)

// BatchProcessorConfig 定义批处理器的配置选项
type BatchProcessorConfig struct {
	BatchSize     int           // 批处理大小
	FlushInterval time.Duration // 刷新间隔
	BufferSize    int           // 缓冲区大小
}

// DefaultConfig 返回默认配置
func DefaultConfig() BatchProcessorConfig {
	return BatchProcessorConfig{
		BatchSize:     100,
		FlushInterval: time.Second,
		BufferSize:    1000,
	}
}

type BatchProcessor[T any] struct {
	batchSize int
	processor func([]T)
	items     chan T
	wg        sync.WaitGroup
	stopChan  chan struct{}
}

// NewBatchProcessor 创建一个新的批处理器
func NewBatchProcessor[T any](batchSize int, processor func([]T)) *BatchProcessor[T] {
	return &BatchProcessor[T]{
		batchSize: batchSize,
		processor: processor,
		items:     make(chan T, batchSize*2),
		stopChan:  make(chan struct{}),
	}
}

// Start 启动批处理器
func (bp *BatchProcessor[T]) Start() {
	bp.wg.Add(1)
	go func() {
		defer bp.wg.Done()
		batch := make([]T, 0, bp.batchSize)
		ticker := time.NewTicker(200 * time.Millisecond) // 添加定时刷新功能
		defer ticker.Stop()

		for {
			select {
			case item := <-bp.items:
				batch = append(batch, item)
				if len(batch) >= bp.batchSize {
					bp.processBatch(batch)
					batch = make([]T, 0, bp.batchSize)
				}
			case <-ticker.C:
				if len(batch) > 0 {
					bp.processBatch(batch)
					batch = make([]T, 0, bp.batchSize)
				}
			case <-bp.stopChan:
				if len(batch) > 0 {
					bp.processBatch(batch)
				}
				return
			}
		}
	}()
}

// Stop 停止批处理器
func (bp *BatchProcessor[T]) Stop() {
	close(bp.stopChan)
	bp.wg.Wait()
}

// Add 添加一个项目到批处理器
func (bp *BatchProcessor[T]) Add(item T) {
	bp.items <- item
}

// processBatch 处理一批数据
func (bp *BatchProcessor[T]) processBatch(batch []T) {
	bp.processor(batch)
}

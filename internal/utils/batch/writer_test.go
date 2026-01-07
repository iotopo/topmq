package batch

import (
	"testing"
	"time"
)

func TestBatchWriter(t *testing.T) {
	// 测试场景1: 测试关闭后追加数据
	// 预期：关闭后 Append 应该返回 false
	t.Run("AppendAfterClose", func(t *testing.T) {
		writer := NewBatchWriter(10, 10*time.Second, func(datas []int) {})
		writer.Close()
		if writer.Append(1) != false {
			t.Fatalf("关闭后 Append 应该返回 false")
		}
	})

	// 测试场景2: 测试达到批次大小后立即写入
	// 预期：当数据量达到 batchSize 时，应该立即触发写入
	t.Run("WriteOnBatchSizeReached", func(t *testing.T) {
		writeCount := 0
		writer := NewBatchWriter(2, 10*time.Second, func(datas []int) {
			writeCount++
		})
		defer writer.Close()

		writer.Append(1)
		writer.Append(2)
		time.Sleep(100 * time.Millisecond) // 等待写入完成

		if writeCount != 1 {
			t.Fatalf("达到批次大小时应该触发一次写入，实际写入次数: %d", writeCount)
		}
	})

	// 测试场景3: 测试定时写入
	// 预期：在未达到批次大小的情况下，应该按照指定的时间间隔写入
	t.Run("WriteOnInterval", func(t *testing.T) {
		writeCount := 0
		writer := NewBatchWriter(10, 100*time.Millisecond, func(datas []int) {
			writeCount++
		})
		defer writer.Close()

		writer.Append(1)
		time.Sleep(150 * time.Millisecond) // 等待超过一个时间间隔

		if writeCount != 1 {
			t.Fatalf("时间间隔到达时应该触发一次写入，实际写入次数: %d", writeCount)
		}
	})

	// 测试场景4: 测试无限缓冲区模式
	// 预期：当 AdditionalBufferSize < 0 时，Append 不应该阻塞
	t.Run("UnboundedBuffer", func(t *testing.T) {
		writeCount := 0
		writer := NewWriter(WriterOptions[int]{
			WriteFunc:            func(datas []int) { writeCount++ },
			BatchSize:            2,
			BatchInterval:        10 * time.Second,
			AdditionalBufferSize: -1,
		})
		defer writer.Close()

		// 快速写入大量数据，不应该阻塞
		for i := 0; i < 1000; i++ {
			if !writer.Append(i) {
				t.Fatalf("无限缓冲区模式下 Append 不应该返回 false")
			}
		}
	})
}

package utils

import (
	"context"
	"sync"
	"time"
)

// LockManager 动态锁管理器
type LockManager struct {
	locks sync.Map      // 存储锁对象和最后使用时间
	stop  chan struct{} // 用于停止自动清理
}

// lockInfo 存储锁对象和最后使用时间
type lockInfo struct {
	lock       *sync.Mutex
	lastAccess time.Time
}

// Lock 获取指定key的锁
func (lm *LockManager) Lock(key string) {
	// 获取或创建锁
	info, _ := lm.locks.LoadOrStore(key, &lockInfo{
		lock:       &sync.Mutex{},
		lastAccess: time.Now(),
	})
	// 更新最后访问时间
	info.(*lockInfo).lastAccess = time.Now()
	// 加锁
	info.(*lockInfo).lock.Lock()
}

// Unlock 释放指定key的锁
func (lm *LockManager) Unlock(key string) {
	// 获取锁对象
	if info, ok := lm.locks.Load(key); ok {
		// 更新最后访问时间
		info.(*lockInfo).lastAccess = time.Now()
		// 解锁
		info.(*lockInfo).lock.Unlock()
	}
}

// TryLock 尝试获取锁，如果获取失败则返回false
func (lm *LockManager) TryLock(key string) bool {
	// 获取或创建锁
	info, _ := lm.locks.LoadOrStore(key, &lockInfo{
		lock:       &sync.Mutex{},
		lastAccess: time.Now(),
	})
	// 更新最后访问时间
	info.(*lockInfo).lastAccess = time.Now()
	// 尝试获取锁
	return info.(*lockInfo).lock.TryLock()
}

// Cleanup 清理长时间未使用的锁
func (lm *LockManager) Cleanup(timeout time.Duration) {
	now := time.Now()
	lm.locks.Range(func(key, value interface{}) bool {
		info := value.(*lockInfo)
		// 检查是否超时
		if now.Sub(info.lastAccess) > timeout {
			// 尝试获取锁，如果成功则说明锁未被使用
			if info.lock.TryLock() {
				// 如果锁未被使用，则删除它
				lm.locks.Delete(key)
				info.lock.Unlock()
			}
		}
		return true
	})
}

// StartAutoCleanup 启动自动清理
func (lm *LockManager) StartAutoCleanup(ctx context.Context, interval, timeout time.Duration) {
	// 如果已经存在清理goroutine，先停止它
	if lm.stop != nil {
		close(lm.stop)
	}
	lm.stop = make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lm.Cleanup(timeout)
			case <-lm.stop:
				return
			}
		}
	}()
}

// StopAutoCleanup 停止自动清理
func (lm *LockManager) StopAutoCleanup() {
	if lm.stop != nil {
		close(lm.stop)
		lm.stop = nil
	}
}

// NewLockManager 创建一个新的锁管理器
func NewLockManager() *LockManager {
	return &LockManager{
		stop: nil,
	}
}

// GlobalLockManager 全局锁管理器实例
var GlobalLockManager = NewLockManager()

// RemoveLock 手动移除指定的锁
// 如果锁正在使用中，会等待锁释放后再删除
// 返回值表示锁是否存在并被删除
func (lm *LockManager) RemoveLock(key string) bool {
	if info, ok := lm.locks.Load(key); ok {
		lockInfo := info.(*lockInfo)
		// 获取锁（如果锁正在使用中会阻塞等待）
		lockInfo.lock.Lock()
		// 删除锁
		lm.locks.Delete(key)
		lockInfo.lock.Unlock()
		return true
	}
	return false
}

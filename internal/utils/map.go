package utils

import (
	"sync"
	"iter"
)

type Map[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// Set 设置键值对
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.m == nil {
		m.m = make(map[K]V)
	}
	m.m[key] = value
}

// Has 检查键是否存在
func (m *Map[K, V]) Has(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.m == nil {
		return false
	}
	_, exists := m.m[key]
	return exists
}

// Get 获取键对应的值
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.m == nil {
		var zero V
		return zero, false
	}
	value, exists := m.m[key]
	return value, exists
}

// SetIfAbsent 仅在键不存在时设置值，返回是否设置成功
func (m *Map[K, V]) SetIfAbsent(key K, value V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.m == nil {
		m.m = make(map[K]V)
	}
	
	if _, exists := m.m[key]; exists {
		return false
	}
	
	m.m[key] = value
	return true
}

// Clear 清空 map
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.m = make(map[K]V)
}

// Size 返回 map 中的元素数量
func (m *Map[K, V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.m == nil {
		return 0
	}
	return len(m.m)
}

// Keys 返回 map 中的所有键
func (m *Map[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.m == nil {
		return []K{}
	}
	
	keys := make([]K, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	return keys
}

// Delete 删除指定的键
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.m != nil {
		delete(m.m, key)
	}
}

// DeleteAndGet 删除并返回指定键的值
func (m *Map[K, V]) DeleteAndGet(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.m == nil {
		var zero V
		return zero, false
	}
	
	value, exists := m.m[key]
	if exists {
		delete(m.m, key)
	}
	return value, exists
}

// Range 返回一个支持 for range 语法的迭代器
func (m *Map[K, V]) Range() iter.Seq2[K, V] {
	m.mu.RLock()
	if m.m == nil {
		m.mu.RUnlock()
		return func(yield func(K, V) bool) {}
	}
	
	// 创建快照
	snapshot := make(map[K]V, len(m.m))
	for k, v := range m.m {
		snapshot[k] = v
	}
	m.mu.RUnlock()
	
	return func(yield func(K, V) bool) {
		for k, v := range snapshot {
			if !yield(k, v) {
				return
			}
		}
	}
}

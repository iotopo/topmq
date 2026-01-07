package cmap

import (
	"iter"
	"sync"
)

// Map 是一个线程安全的泛型 Map
// K 是键的类型，必须是可比较的类型（comparable）
// V 是值的类型，可以是任意类型（any）
type Map[K comparable, V any] struct {
	internal sync.Map
}

// New 创建一个新的泛型 Map
// 返回一个新的空 Map 实例
func New[K comparable, V any]() Map[K, V] {
	return Map[K, V]{}
}

// Store 存储键值对
// 参数：
//   - key: 要存储的键
//   - value: 要存储的值
//
// 如果键已存在，会更新其值
func (m *Map[K, V]) Store(key K, value V) {
	m.internal.Store(key, value)
}

// StoreIfAbsent 仅在键不存在时存储值
// 参数：
//   - key: 要存储的键
//   - value: 要存储的值
//
// 返回值：
//   - V: 实际的值（如果键存在则是已存在的值，如果键不存在则是新存储的值）
//   - bool: 是否成功存储了新值（true 表示存储了新值，false 表示键已存在）
func (m *Map[K, V]) StoreIfAbsent(key K, value V) (V, bool) {
	actual, loaded := m.internal.LoadOrStore(key, value)
	return actual.(V), !loaded
}

// Has 判断键是否存在
// 参数：
//   - key: 要查找的键
//
// 返回值：
//   - bool: 是否存在键（true 表示存在，false 表示不存在）
func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.internal.Load(key)
	return ok
}

// Load 加载键对应的值
// 参数：
//   - key: 要查找的键
//
// 返回值：
//   - V: 找到的值，如果键不存在则返回零值
//   - bool: 是否找到键（true 表示找到，false 表示未找到）
func (m *Map[K, V]) Load(key K) (V, bool) {
	value, ok := m.internal.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return value.(V), true
}

// LoadOrStore 加载键对应的值，如果不存在则存储默认值
// 参数：
//   - key: 要查找的键
//   - value: 如果键不存在时要存储的值
//
// 返回值：
//   - V: 实际的值（如果键存在则是已存在的值，如果键不存在则是新存储的值）
//   - bool: 是否加载了已存在的值（true 表示加载了已存在的值，false 表示存储了新值）
func (m *Map[K, V]) LoadOrStore(key K, value V) (V, bool) {
	actual, loaded := m.internal.LoadOrStore(key, value)
	return actual.(V), loaded
}

// Delete 删除键值对
// 参数：
//   - key: 要删除的键
//
// 如果键不存在，此操作不会产生任何效果
func (m *Map[K, V]) Delete(key K) {
	m.internal.Delete(key)
}

// Range 遍历所有键值对
// 参数：
//   - f: 遍历函数，接收键和值作为参数，返回是否继续遍历
//
// 遍历函数返回 false 可以提前结束遍历
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.internal.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

// Len 返回 Map 中的键值对数量
// 返回值：
//   - int: Map 中的键值对数量
//
// 注意：此方法需要遍历整个 Map，性能较差
func (m *Map[K, V]) Len() int {
	count := 0
	m.Range(func(key K, value V) bool {
		count++
		return true
	})
	return count
}

// Clear 清空 Map
// 删除 Map 中的所有键值对
// 注意：此方法需要遍历整个 Map，性能较差
func (m *Map[K, V]) Clear() {
	m.Range(func(key K, value V) bool {
		m.Delete(key)
		return true
	})
}

// Pop 根据指定的 key 从 Map 中移除并返回对应的值，同 LoadAndDelete
// 参数：
//   - key: 要移除的键
//
// 返回值：
//   - V: 被移除的值，如果键不存在则返回零值
//   - bool: 是否成功移除（true 表示成功移除，false 表示键不存在）
func (m *Map[K, V]) Pop(key K) (V, bool) {
	return m.LoadAndDelete(key)
}

// LoadAndDelete 根据指定的 key 从 Map 中移除并返回对应的值
// 参数：
//   - key: 要移除的键
//
// 返回值：
//   - V: 被移除的值，如果键不存在则返回零值
//   - bool: 是否成功移除（true 表示成功移除，false 表示键不存在）
func (m *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	actual, loaded := m.internal.LoadAndDelete(key)
	if !loaded {
		var zero V
		return zero, false
	}
	return actual.(V), loaded
}

// Iter 返回一个迭代器，支持 for range 语法
// 返回值：
//   - iter.Seq2[K, V]: 用于遍历键值对的迭代器
//
// 使用示例：
//
//	m := cmap.New[string, int]()
//	m.Store("key1", 1)
//	m.Store("key2", 2)
//
//	// 基本用法
//	for key, value := range m.Iter() {
//	    fmt.Printf("key: %s, value: %d\n", key, value)
//	}
//
//	// 提前退出
//	for key, value := range m.Iter() {
//	    if someCondition {
//	        break
//	    }
//	    // 处理键值对
//	}
//
//	// 收集所有键
//	keys := make([]string, 0)
//	for key, _ := range m.Iter() {
//	    keys = append(keys, key)
//	}
//
//	// 收集所有值
//	values := make([]int, 0)
//	for _, value := range m.Iter() {
//	    values = append(values, value)
//	}
func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(func(key K, value V) bool {
			return yield(key, value)
		})
	}
}

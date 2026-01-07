package cmap

import (
	"sync"
	"testing"
)

func TestMap_BasicOperations(t *testing.T) {
	m := New[string, int]()

	// 测试 Store 和 Load
	m.Store("key1", 1)
	value, ok := m.Load("key1")
	if !ok || value != 1 {
		t.Errorf("Load failed: expected 1, got %v, ok: %v", value, ok)
	}

	// 测试 Load 不存在的键
	value, ok = m.Load("nonexistent")
	if ok {
		t.Error("Load nonexistent key should return false")
	}

	// 测试 Delete
	m.Delete("key1")
	value, ok = m.Load("key1")
	if ok {
		t.Error("Deleted key should not exist")
	}

	// 测试 LoadOrStore
	value, loaded := m.LoadOrStore("key2", 2)
	if loaded {
		t.Error("LoadOrStore should not report loaded for new key")
	}
	if value != 2 {
		t.Errorf("LoadOrStore failed: expected 2, got %v", value)
	}

	// 测试 LoadAndDelete
	value, ok = m.LoadAndDelete("key2")
	if !ok || value != 2 {
		t.Errorf("LoadAndDelete failed: expected 2, got %v, ok: %v", value, ok)
	}
	value, ok = m.Load("key2")
	if ok {
		t.Error("LoadAndDelete key should not exist")
	}
}

func TestMap_Range(t *testing.T) {
	m := New[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)
	m.Store("key3", 3)

	// 测试 Range
	keys := make(map[string]bool)
	values := make(map[int]bool)
	m.Range(func(key string, value int) bool {
		keys[key] = true
		values[value] = true
		return true
	})

	if len(keys) != 3 || len(values) != 3 {
		t.Error("Range did not iterate over all elements")
	}

	// 测试 Range 提前返回
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return count < 2 // 只遍历前两个元素
	})
	if count != 2 {
		t.Errorf("Range should stop after 2 elements, got %d", count)
	}
}

func TestMap_Concurrent(t *testing.T) {
	m := New[string, int]()
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Store("key"+string(rune(i)), i)
		}(i)
	}
	wg.Wait()

	// 并发读取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value, ok := m.Load("key" + string(rune(i)))
			if !ok {
				t.Errorf("Failed to load key%d", i)
			}
			if value != i {
				t.Errorf("Expected %d, got %d", i, value)
			}
		}(i)
	}
	wg.Wait()

	// 并发删除
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Delete("key" + string(rune(i)))
		}(i)
	}
	wg.Wait()

	// 验证所有元素都被删除
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("Expected empty map, got %d elements", count)
	}
}

func TestMap_Len(t *testing.T) {
	m := New[string, int]()
	if m.Len() != 0 {
		t.Error("Empty map should have length 0")
	}

	m.Store("key1", 1)
	m.Store("key2", 2)
	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}

	m.Delete("key1")
	if m.Len() != 1 {
		t.Errorf("Expected length 1, got %d", m.Len())
	}
}

func TestMap_Clear(t *testing.T) {
	m := New[string, int]()
	m.Store("key1", 1)
	m.Store("key2", 2)

	m.Clear()
	if m.Len() != 0 {
		t.Error("Map should be empty after Clear")
	}

	_, ok := m.Load("key1")
	if ok {
		t.Error("Map should not contain elements after Clear")
	}
}

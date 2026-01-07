package utils

import (
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("Basic Operations", func(t *testing.T) {
		var m Map[string, int]

		// Set 和 Get
		m.Set("a", 1)
		if val, exists := m.Get("a"); !exists || val != 1 {
			t.Errorf("expected (1, true), got (%v, %v)", val, exists)
		}

		// Has
		if !m.Has("a") {
			t.Error("expected key 'a' to exist")
		}
		if m.Has("non-existent") {
			t.Error("expected key 'non-existent' to not exist")
		}

		// SetIfAbsent
		if success := m.SetIfAbsent("a", 2); success {
			t.Error("SetIfAbsent should return false for existing key")
		}
		if success := m.SetIfAbsent("b", 2); !success {
			t.Error("SetIfAbsent should return true for new key")
		}
		if val, _ := m.Get("a"); val != 1 {
			t.Errorf("expected value 1, got %v", val)
		}
		if val, _ := m.Get("b"); val != 2 {
			t.Errorf("expected value 2, got %v", val)
		}

		// Delete
		m.Delete("a")
		if m.Has("a") {
			t.Error("key 'a' should be deleted")
		}

		// DeleteAndGet
		m.Set("c", 3)
		if val, exists := m.DeleteAndGet("c"); !exists || val != 3 {
			t.Errorf("expected (3, true), got (%v, %v)", val, exists)
		}
		if m.Has("c") {
			t.Error("key 'c' should be deleted")
		}

		// Size
		m.Clear()
		if size := m.Size(); size != 0 {
			t.Errorf("expected size 0, got %v", size)
		}
		m.Set("a", 1)
		m.Set("b", 2)
		if size := m.Size(); size != 2 {
			t.Errorf("expected size 2, got %v", size)
		}

		// Keys
		keys := m.Keys()
		if len(keys) != 2 {
			t.Errorf("expected 2 keys, got %v", len(keys))
		}
		keyMap := make(map[string]bool)
		for _, k := range keys {
			keyMap[k] = true
		}
		if !keyMap["a"] || !keyMap["b"] {
			t.Error("missing expected keys")
		}
	})

	t.Run("Iterator Operations", func(t *testing.T) {
		var m Map[string, int]
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		// 基本迭代
		count := 0
		sum := 0
		for k, v := range m.Range() {
			count++
			sum += v
			if v != 1 && v != 2 && v != 3 {
				t.Errorf("unexpected value %v for key %v", v, k)
			}
		}
		if count != 3 {
			t.Errorf("expected 3 iterations, got %v", count)
		}
		if sum != 6 {
			t.Errorf("expected sum 6, got %v", sum)
		}

		// 迭代过程中的修改
		for k, v := range m.Range() {
			m.Set("new-key", 100)
			m.Delete(k)
			if v != 1 && v != 2 && v != 3 {
				t.Errorf("unexpected value %v for key %v", v, k)
			}
		}
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		var m Map[int, int]
		const goroutines = 10
		const operationsPerGoroutine = 1000

		var wg sync.WaitGroup
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func(base int) {
				defer wg.Done()
				for j := 0; j < operationsPerGoroutine; j++ {
					key := base*operationsPerGoroutine + j
					m.Set(key, key)
					if val, exists := m.Get(key); !exists || val != key {
						t.Errorf("expected (%v, true), got (%v, %v)", key, val, exists)
					}
				}
			}(i)
		}
		wg.Wait()

		if size := m.Size(); size != goroutines*operationsPerGoroutine {
			t.Errorf("expected size %v, got %v", goroutines*operationsPerGoroutine, size)
		}
	})

	t.Run("Nil Map Operations", func(t *testing.T) {
		var m Map[string, int]

		if m.Has("key") {
			t.Error("Has should return false for nil map")
		}

		if val, exists := m.Get("key"); exists || val != 0 {
			t.Errorf("Get should return (0, false) for nil map, got (%v, %v)", val, exists)
		}

		if size := m.Size(); size != 0 {
			t.Errorf("Size should return 0 for nil map, got %v", size)
		}

		if keys := m.Keys(); len(keys) != 0 {
			t.Errorf("Keys should return empty slice for nil map, got %v", keys)
		}

		count := 0
		for _, _ = range m.Range() {
			count++
		}
		if count != 0 {
			t.Errorf("Range should not iterate over nil map, got %v iterations", count)
		}

		// 自动初始化
		m.Set("key", 1)
		if !m.Has("key") {
			t.Error("Map should be automatically initialized after Set")
		}
	})
}

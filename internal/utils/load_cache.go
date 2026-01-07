package utils

import "sync"

type LoadCacheItem[V any] struct {
	Error error
	Value V
	Done  chan struct{}
}

type LoadCache[K comparable, V any] struct {
	inner sync.Map
}

// 如果 key 对应的值没有加载过，调用一次 loader 进行加载
// 无论调用成功与否，后续以相同的 key 都会返回此调用结果
// 并发调用时，相同 key 的 loader 保证只调用一次
func (c *LoadCache[K, V]) Load(key K, loader func(K) (V, error)) (V, bool, error) {
	done := make(chan struct{})
	defer close(done)

	value, loaded := c.inner.LoadOrStore(key, &LoadCacheItem[V]{Done: done})
	cached := value.(*LoadCacheItem[V])
	if loaded {
		<-cached.Done
	} else {
		cached.Value, cached.Error = loader(key)
	}
	return cached.Value, loaded, cached.Error
}

package lru

import (
	"container/list"
	"sync"
	"time"
)

type CacheEntry[K comparable, V any] struct {
	Key           K
	Value         V
	LastVisitTime time.Time
}

// Cache 是一个支持 TTL 的 LRU(Least Recently Used) 缓存
// 缓存满时会删除最长时间没有使用的元素
// 和其他支持 TTL 的 LRU 缓存不同的是，当访问元素时，其 TTL 会重置
type Cache[K comparable, V any] struct {
	mu         sync.Mutex
	maxEntries int
	ttl        time.Duration
	ll         *list.List
	cache      map[K]*list.Element
	onEvict    func(K, V)
	stop       chan struct{}
}

// Option 定义缓存配置选项类型
type Option[K comparable, V any] func(*Cache[K, V])

// WithMaxEntries 设置缓存最大条目数
func WithMaxEntries[K comparable, V any](maxEntries int) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.maxEntries = maxEntries
	}
}

// WithTTL 设置缓存条目的过期时间
func WithTTL[K comparable, V any](ttl time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.ttl = ttl
	}
}

// WithOnEvict 设置缓存条目被移除时的回调函数
func WithOnEvict[K comparable, V any](onEvict func(K, V)) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.onEvict = onEvict
	}
}

func (c *Cache[K, V]) init() {
	c.cache = make(map[K]*list.Element)
	c.ll = list.New()
}

func (c *Cache[K, V]) removeOldest() {
	if elem := c.ll.Back(); elem != nil {
		c.ll.Remove(elem)

		entry := elem.Value.(*CacheEntry[K, V])
		if c.onEvict != nil {
			c.onEvict(entry.Key, entry.Value)
		}
	}
}

func (c *Cache[K, V]) Store(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.init()
	}

	if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)

		entry := elem.Value.(*CacheEntry[K, V])
		if c.onEvict != nil {
			c.onEvict(key, entry.Value)
		}
		entry.Value = value
		if c.ttl > 0 {
			entry.LastVisitTime = time.Now()
		}
		return
	}
	elem := c.ll.PushFront(&CacheEntry[K, V]{Key: key, Value: value, LastVisitTime: time.Now()})
	c.cache[key] = elem
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

func (c *Cache[K, V]) LoadOrStore(key K, valueFunc func() (V, error)) (value V, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.init()
	} else if elem, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elem)

		entry := elem.Value.(*CacheEntry[K, V])
		if c.ttl > 0 {
			entry.LastVisitTime = time.Now()
		}
		return entry.Value, nil
	}

	value, err = valueFunc()
	if err != nil {
		return
	}

	elem := c.ll.PushFront(&CacheEntry[K, V]{Key: key, Value: value, LastVisitTime: time.Now()})
	c.cache[key] = elem
	if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
	return
}

func (c *Cache[K, V]) Load(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}
	if elem, hit := c.cache[key]; hit {
		c.ll.MoveToFront(elem)

		entry := elem.Value.(*CacheEntry[K, V])
		if c.ttl > 0 {
			entry.LastVisitTime = time.Now()
		}
		return entry.Value, true
	}
	return
}

func (c *Cache[K, V]) LoadAndDelete(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}

	if elem, hit := c.cache[key]; hit {
		entry := elem.Value.(*CacheEntry[K, V])
		if c.onEvict != nil {
			c.onEvict(entry.Key, entry.Value)
		}
		c.ll.Remove(elem)
		delete(c.cache, entry.Key)
		return entry.Value, true
	}
	return
}

func (c *Cache[K, V]) doClear() {
	if c.cache == nil {
		return
	}

	if c.onEvict != nil {
		for _, e := range c.cache {
			entry := e.Value.(*CacheEntry[K, V])
			c.onEvict(entry.Key, entry.Value)
		}
	}
	c.ll = nil
	c.cache = nil
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.doClear()
}

func (c *Cache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *Cache[K, V]) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stop != nil {
		close(c.stop)
		c.stop = nil
	}
	c.doClear()
}

func (c *Cache[K, V]) doExpire() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}

	now := time.Now()

	for {
		if elem := c.ll.Back(); elem != nil {
			entry := elem.Value.(*CacheEntry[K, V])
			if now.Sub(entry.LastVisitTime) > c.ttl {
				if c.onEvict != nil {
					c.onEvict(entry.Key, entry.Value)
				}
				c.ll.Remove(elem)
				delete(c.cache, entry.Key)
			} else {
				return
			}
		} else {
			return
		}
	}
}

// NewCache 创建一个新的 LRU 缓存实例
// 使用 WithOption 模式配置缓存的各种参数
func NewCache[K comparable, V any](options ...Option[K, V]) *Cache[K, V] {
	cache := &Cache[K, V]{
		mu: sync.Mutex{},
		// 默认值设置
		maxEntries: 0, // 0 表示不限制
		ttl:        0, // 0 表示永不过期
	}

	// 应用所有配置选项
	for _, opt := range options {
		opt(cache)
	}

	// 初始化缓存结构
	cache.init()

	// 如果设置了 TTL，则启动过期检查协程
	if cache.ttl > 0 {
		stop := make(chan struct{})
		cache.stop = stop

		go func() {
			ticker := time.NewTicker(cache.ttl)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					cache.doExpire()
				case <-stop:
					return
				}
			}
		}()
	}

	return cache
}

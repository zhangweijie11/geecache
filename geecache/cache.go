package geecache

import (
	"geecache/geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex // 互斥锁
	lru        *lru.Cache // 缓存数据的操作
	cacheBytes int64      // 缓存占用内存值
}

// 添加数据
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil { // 延迟初始化
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// 查询缓存
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

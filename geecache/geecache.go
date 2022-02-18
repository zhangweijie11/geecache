package geecache

import (
	"fmt"
	"log"
	"sync"
)

// 查询数据处理器
type Getter interface {
	Get(key string) ([]byte, error)
}

// 处理器函数
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function, 缓存未命中时进行源数据回调
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 缓存组命名空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// Group 构造函数
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// 获取缓存组
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// 获取缓存值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	v, ok := g.mainCache.get(key)
	if ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}

// 不存在缓存时回调获取源数据
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// 获取源数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 缓存源数据
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

package lru

import "container/list"

type Cache struct {
	maxBytes  int64                         // 缓存内存最大值
	nbytes    int64                         // 当前使用的内存值，内存使用是键值都会使用
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 缓存键值映射关系
	OnEvictes func(key string, value Value) //  回调函数，当缓存不存在或者其他情况下调用
}

// 缓存键值对，模拟双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// 返回值占用的内存大小
type Value interface {
	Len() int
}

// Cache 构造函数
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     map[string]*list.Element{},
		OnEvictes: onEvicted,
	}
}

// 查找功能
func (c *Cache) Get(key string) (value Value, ok bool) {
	/*
		如果获取到数据曾返回数据值，否则返回空
	*/
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) // 将链表中的节点移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除功能,按照 LRU 算法的淘汰策略，直接淘汰最近最少访问的节点即队首
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() // 获取队首节点
	if ele != nil {
		c.ll.Remove(ele) // 双向链表删除队首节点
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // 删除映键值对映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) // 更新已使用内存值
		if c.OnEvictes != nil {                                // 调用回调函数
			c.OnEvictes(kv.key, kv.value)
		}
	}
}

// 新增/修改
func (c *Cache) Add(key string, value Value) {
	// 当节点不存在时添加节点
	if ele, ok := c.cache[key]; !ok {
		ele := c.ll.PushFront(&entry{key: key, value: value}) // 双向链表添加节点
		c.cache[key] = ele                                    // 增加键值对映射关系
		c.nbytes += int64(len(key)) + int64(value.Len())      // 更新缓存已使用的内存值
	} else {
		c.ll.MoveToFront(ele) // 将节点移至队首
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) // 更新缓存已使用的内存值
		kv.value = value
	}
	// 判断当前已使用内存值是否超过最大值，如果超过了则循环删除至最大值之下
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 统计缓存数据数量
func (c *Cache) Len() int {
	return c.ll.Len()
}

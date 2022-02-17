package lru

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

// 测试获取数据
func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	// 如果无法获取 key1 的值或者值不是 1234 则失败
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatal("cache hit key1=1234 failed")
	}
	// 如果能获取 key2 的值则失败，因为缓存中不可能存在 key2
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

// 测试内存超过限定值，触发节点清除功能
func TestCache_RemoveOldest(t *testing.T) {
	// 添加三条数据，但是设定了内存最大值，超过内存最大值最先添加的缓存数据会被删除
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	caps := len(k1 + k2 + v1 + v2)
	lru := New(int64(caps), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatal("RemoveOldest key1 failed")
	}
}

// 测试回调函数
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	// 将回调函数变成一个不定长链表
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	// 设定了内存最大值，无法容纳全部缓存数据
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	// 预期被删除数据为 key1 和 k2
	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

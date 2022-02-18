package geecache

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetGroup(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("source", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			// 如果命中缓存，则命中数加一，反之则命中数值为零， 返回错误信息
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cacahe %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("this value of unknow should be empty,but %s got", view)
	}
}

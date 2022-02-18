package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	self     string // 本地地址 ip：port
	basePath string // URL 前缀
}

// HTTPPool 构造函数
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 日志记录功能
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 验证请求的路径前缀和默认的前缀是否一致
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	// 记录日志，当前日志没有级别
	p.Log("%s %s", r.Method, r.URL.Path)
	// 限制 URL 格式 /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusNotFound)
		return
	}

	groupName := parts[0] // 缓存组名称
	key := parts[1]       // 缓存键
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group", http.StatusNotFound)
		return
	}
	// 获取缓存数据
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// 设置数据返回类型为字节
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(view.ByteSlice()); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

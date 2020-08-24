package gocache

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

const defaultBasePath = "/_gocache/"

// HTTPPool contains a peer's name and path
type HTTPPool struct {
	self     string
	basePath string
}

// NewHTTPPool creates an instace of http server which is used to
// communicate with the other peers
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log logs info with server name
func (pool *HTTPPool) Log(format string, value ...interface{}) {
	log.Printf("[gocache Server %s] %s", pool.self, fmt.Sprintf(format, value...))
}

// ServeHTTP serves HTTP connections
func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pattern := path.Join(pool.basePath, "*/[a-zA-Z0-9]*")
	// check pattern matched
	// /_gocache/:group/:id
	if matched, err := path.Match(pattern, r.URL.Path); !matched || err != nil {
		http.Error(w, "path: "+pattern+" expexted, but get "+r.URL.Path, http.StatusNotFound)
		return
	}
	pool.Log("%s %s", r.Method, r.URL.Path)
	groupName := strings.TrimPrefix(path.Dir(r.URL.Path), pool.basePath)
	key := path.Base(r.URL.Path)
	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, "group: "+groupName+" not exists!", http.StatusNotFound)
	}
	bv, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bv.ByteSlice())
}

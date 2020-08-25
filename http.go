package gocache

import (
	"fmt"
	"gocache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_gocache/"
	defaultReplicas = 50
)

// HTTPPool is a http server contains a peer's name and path
type HTTPPool struct {
	// this peer's host name and port number
	// eg: https//go.sino.moe:8000
	self string
	// api prefix
	// eg: https//go.sino.moe:8000{basePath}{group}/{key}
	basePath string

	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
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

// SetPeers sets pool's list of peers
func (pool *HTTPPool) SetPeers(peers ...string) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.peers = consistenthash.New(defaultReplicas, nil)
	pool.peers.Add(peers...)
	// dalay initialization
	if pool.httpGetters == nil {
		pool.httpGetters = make(map[string]*httpGetter)
	}
	for _, peer := range peers {
		pool.httpGetters[peer] = &httpGetter{
			baseURL: peer + pool.basePath,
		}
	}
}

// PickPeer picks a peer according to the specified key
func (pool *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if peer := pool.peers.Get(key); peer != "" && peer != pool.self {
		pool.Log("pick peer: %s", peer)
		return pool.httpGetters[key], false
	}
	return nil, false
}

// interface assertion
var (
	_ PeerGetter = (*httpGetter)(nil)
	_ PeerPicker = (*HTTPPool)(nil)
)

// httpGetter is a http client to get remote peers' data
type httpGetter struct {
	// remote peer's api base path
	// by default is defaultBasePath
	baseURL string
}

// Get uses an http client to get remote peers data
func (h *httpGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote peer returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %v", err)
	}

	return bytes, nil
}

package octopus

import (
	"sync"
)

type routes struct {
	sync.RWMutex
	data map[string]*route
}

type route struct {
	sync.RWMutex
	data             map[string][]HandlerFunc
	a                *app
	globalMiddleware []HandlerFunc
	path             string
}

func (rs *routes) add(path string, method string, handler ...HandlerFunc) {
	rs.Lock()
	defer rs.Unlock()
	if rs.data == nil {
		rs.data = make(map[string]*route)
	}
	if rs.data[path] == nil {
		rs.data[path] = &route{data: make(map[string][]HandlerFunc)}
	}
	rs.data[path].data[method] = handler
}

func (rs *routes) get(path string, method string) []HandlerFunc {
	rs.RLock()
	defer rs.RUnlock()
	if rs.data == nil {
		return nil
	}
	if rs.data[path].data == nil {
		return nil
	}
	return rs.data[path].data[method]
}

func (rs *routes) rrange(f func(key string, value *route) bool) {
	rs.RLock()
	defer rs.RUnlock()
	for k, v := range rs.data {
		if !f(k, v) {
			break
		}
	}
}

func (r *route) methodExists(method string) ([]HandlerFunc, bool) {
	r.RLock()
	defer r.RUnlock()
	hs, exists := r.data[method]
	return hs, exists
}

func (r *route) Get(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.Get(r.path+path, handlers...)
}

func (r *route) DELETE(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.Post(r.path+path, handlers...)
}

func (r *route) PUT(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.PUT(r.path+path, handlers...)
}

func (r *route) Post(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.Post(r.path+path, handlers...)
}

func (r *route) PATCH(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.PATCH(r.path+path, handlers...)
}

func (r *route) OPTIONS(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.OPTIONS(r.path+path, handlers...)
}

func (r *route) HEAD(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.HEAD(r.path+path, handlers...)
}

func (r *route) Any(path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.Any(r.path+path, handlers...)
}

func (r *route) Method(method string, path string, handlers ...HandlerFunc) {
	handlers = append(r.globalMiddleware, handlers...)
	r.a.Method(method, r.path+path, handlers...)
}

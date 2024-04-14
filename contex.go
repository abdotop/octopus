package octopus

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
)

type Ctx struct {
	// sync.RWMutex
	handlers []HandlerFunc
	index    int
	Values   *value
	Context  context.Context
}

func NewCtx() *Ctx {
	return &Ctx{
		handlers: []HandlerFunc{},
		index:    0,
		Values:   new(value),
		Context:  context.Background(),
	}
}

type Map = map[string]interface{}

func (c *Ctx) BodyParser(out interface{}) error {
	// c.RLock()
	// defer c.RUnlock()
	r, ok := c.Values.Get("request")
	if ok {
		r := r.(*http.Request)
		return json.NewDecoder(r.Body).Decode(&out)
	}
	return errors.New("request not found in context values")
}

// Get returns the value of the key in the context header
func (c *Ctx) Get(key string) string {
	// c.RLock()
	// defer c.RUnlock()
	r, ok := c.Values.Get("request")
	if ok {
		r := r.(*http.Request)
		return r.Header.Get(key)
	}
	return ""
}

func (c *Ctx) JSON(data interface{}) error {
	// c.Lock()
	// defer c.Unlock()
	r, ok := c.Values.Get("response")
	if ok {
		r := r.(http.ResponseWriter)
		r.Header().Set("Content-Type", "application/json")
		return json.NewEncoder(r).Encode(data)
	}
	return errors.New("response not found in context values")
}

func (c *Ctx) Next() {
	if c.index < len(c.handlers) {
		handler := c.handlers[c.index]
		c.index++
		handler(c)
	}
}

func (c *Ctx) Query(key string) string {
	// c.RLock()
	// defer c.RUnlock()
	r, ok := c.Values.Get("request")
	if ok {
		r := r.(*http.Request)
		return r.URL.Query().Get(key)
	}
	return ""
}

func (c *Ctx) Render(path string, data interface{}) error {
	// c.Lock()
	// defer c.Unlock()
	r, ok := c.Values.Get("response")
	if ok {
		r := r.(http.ResponseWriter)
		tp, err := template.ParseFiles(path)
		if err != nil {
			return err
		}
		return tp.Execute(r, data)
	}
	return errors.New("response not found in context values")
}

func (c *Ctx) SendString(code statusCode, s string) error {
	// c.Lock()
	// defer c.Unlock()
	r, ok := c.Values.Get("response")
	if ok {
		c.Status(code)
		r := r.(http.ResponseWriter)
		_, err := r.Write([]byte(s))
		return err
	}
	return errors.New("response not found in context values")
}

func (c *Ctx) Status(code statusCode) *Ctx {
	// c.RLock()
	// defer c.RUnlock()
	r, ok := c.Values.Get("response")
	a, appExist := c.Values.Get("app")
	if ok {
		r := r.(http.ResponseWriter)
		r.WriteHeader(int(code))
		if appExist {
			a := a.(*app)
			a.handleError(code, c)
		} else {
			a := New()
			a.handleError(code, c)
		}
	}
	return c
}

func (c *Ctx) WriteString(s string) error {
	// c.RLock()
	// defer c.RUnlock()
	r, ok := c.Values.Get("response")
	if ok {
		r := r.(http.ResponseWriter)
		_, err := r.Write([]byte(s))
		return err
	}
	return errors.New("response not found in context values")
}

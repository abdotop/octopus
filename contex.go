package octopus

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
)

type Ctx struct {
	Response http.ResponseWriter
	Request  *http.Request
	handlers []Handler
	index    int
	Values   map[any]any
	Context  context.Context
}

type Map = map[string]interface{}

func (c *Ctx) BodyParser(out interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(&out)
}

// Get returns the value of the key in the context header
func (c *Ctx) Get(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Ctx) JSON(data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response).Encode(data)
}

func (c *Ctx) Next() {
	if c.index < len(c.handlers) {
		handler := c.handlers[c.index]
		c.index++
		handler(c)
	}
}

func (c *Ctx) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Ctx) Render(path string, data interface{}) error {
	tp, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	return tp.Execute(c.Response, data)
}

func (c *Ctx) SendString(s string) error {
	_, err := c.Response.Write([]byte(s))
	return err
}

func (c *Ctx) Status(code int) *Ctx {
	c.Response.WriteHeader(code)
	return c
}

func (c *Ctx) WriteString(s string) (int, error) {
	return c.Response.Write([]byte(s))
}

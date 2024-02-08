package octopus

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Ctx struct {
	Response http.ResponseWriter
	Request  *http.Request
	handlers []HandlerFunc
	index    int
	Values   map[any]any
}

func (c *Ctx) BodyParser(out interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(&out)
}

func (c *Ctx) Get(key string, defaultValue ...string) string {
	value := c.Request.Header.Get(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return value
}

func (c *Ctx) JSON(data interface{}) error {
	c.Response.Header().Set("Content-Type", "Application/json")
	return json.NewEncoder(c.Response).Encode(data)
}

func (c *Ctx) Next() {
	if c.index < len(c.handlers) {
		handler := c.handlers[c.index]
		c.index++
		handler(c)
	}
}

func (c *Ctx) Render(path string, data interface{}) error {
	tp, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	return tp.Execute(c.Response, data)
}

func (c *Ctx) Status(code int) *Ctx {
	c.Response.WriteHeader(code)
	return c
}

func (c *Ctx) WriteString(s string) (int, error) {
	return c.Response.Write([]byte(s))
}

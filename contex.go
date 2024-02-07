package octopus

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	handlers       []HandlerFunc
	index          int
}

func (c *Context) BodyParser(out interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(&out)
}

func (c *Context) JSON(data interface{}) error {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.ResponseWriter).Encode(data)
}

func (c *Context) Next() {
	if c.index < len(c.handlers) {
		handler := c.handlers[c.index]
		c.index++
		handler(c)
	}
}

func (c *Context) Render(path string, data interface{}) error {
	tp, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	return tp.Execute(c.ResponseWriter, data)
}

func (c *Context) Status(code int) *Context {
	c.ResponseWriter.WriteHeader(code)
	return c
}

func (c *Context) WriteString(s string) (int, error) {
	return c.ResponseWriter.Write([]byte(s))
}

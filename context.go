package mapes

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	w      http.ResponseWriter
	r      *http.Request
	m      map[string]any
	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r, make(map[string]any), make(map[string]string)}
}

//Response

func (c *Context) Json(code int, value any) error {
	indent, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return err
	}

	c.w.WriteHeader(code)
	_, err = c.w.Write([]byte(indent))
	return err
}

func (c *Context) String(code int, value string) error {
	c.w.WriteHeader(code)
	_, err := c.w.Write([]byte(value))
	return err
}

func (c *Context) None(code int) error {
	c.w.WriteHeader(code)
	return nil
}

//Context

func (c *Context) Get(key string) any {
	return c.m[key]
}

func (c *Context) Set(key string, value any) {
	c.m[key] = value
}

func (c *Context) Context() context.Context {
	return c.r.Context()
}

//Binding

func (c *Context) Bind(dst any) error {
	decoder := json.NewDecoder(c.r.Body)
	return decoder.Decode(dst)
}

//Params

func (c *Context) Param(key string) string {
	return c.params[key]
}

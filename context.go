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

func (c *Context) Json(code int, value any) error {
	c.SetHeader("content-type", "application/json")

	indent, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return err
	}

	c.w.WriteHeader(code)
	_, err = c.w.Write([]byte(indent))
	return err
}

func (c *Context) String(code int, value string) error {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	_, err := c.w.Write([]byte(value))
	return err
}

func (c *Context) Status(code int) {
	c.w.WriteHeader(code)
}

func (c *Context) None(code int) error {
	c.w.WriteHeader(code)
	return nil
}

func (c *Context) Get(key string) any {
	return c.m[key]
}

func (c *Context) Set(key string, value any) {
	c.m[key] = value
}

func (c *Context) Context() context.Context {
	return c.r.Context()
}

func (c *Context) Bind(dst any) error {
	decoder := json.NewDecoder(c.r.Body)
	return decoder.Decode(dst)
}

func (c *Context) Param(key string) string {
	return c.params[key]
}

func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *Context) SetHeader(key, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) AddHeader(key, value string) {
	c.w.Header().Add(key, value)
}

func (c *Context) FormValues(key string) string {
	return c.r.FormValue(key)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.w, cookie)
}

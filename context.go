package mapes

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
	m map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r, make(map[string]interface{})}
}

//Response

func (c *Context) Json(code int, value interface{}) error {
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

//Anu

func (c *Context) Get(key string) interface{} {
	return c.m[key]
}

func (c *Context) Set(key string, value interface{}) {
	c.m[key] = value
}

func (c *Context) Context() context.Context {
	return c.r.Context()
}

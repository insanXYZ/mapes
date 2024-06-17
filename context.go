package mapes

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
	"slices"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	m       map[string]any
	params  map[string]string
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

	c.Writer.WriteHeader(code)
	_, err = c.Writer.Write([]byte(indent))
	return err
}

func (c *Context) String(code int, value string) error {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	_, err := c.Writer.Write([]byte(value))
	return err
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) None(code int) error {
	c.Writer.WriteHeader(code)
	return nil
}

func (c *Context) Get(key string) any {
	return c.m[key]
}

func (c *Context) Set(key string, value any) {
	c.m[key] = value
}

func (c *Context) Context() context.Context {
	return c.Request.Context()
}

func (c *Context) Bind(dst any) error {
	if slices.Contains(c.Request.Header["Content-Type"], "application/json") {
		err := c.bindJson(dst)
		if err != nil {
			return err
		}
	}

	if slices.Contains(c.Request.Header["Content-Type"], "application/xml") {
		err := c.bindXml(dst)
		if err != nil {
			return err
		}
	}

	if slices.Contains(c.Request.Header["Content-Type"], "application/x-www-form-urlencoded") {
		err := c.bindForm(dst)
		if err != nil {
			return err
		}
	}

	c.bindUrl(dst)

	return nil
}

func (c *Context) bindJson(dst any) error {
	decoder := json.NewDecoder(c.Request.Body)
	return decoder.Decode(dst)
}

func (c *Context) bindXml(dst any) error {
	decoder := xml.NewDecoder(c.Request.Body)
	return decoder.Decode(dst)
}

func (c *Context) bindForm(dst any) error {
	err := c.Request.ParseForm()
	if err != nil {
		return err
	}

	v := reflect.ValueOf(dst).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag := field.Tag.Get("form"); tag != "" {
			v.Field(i).SetString(c.Request.FormValue(tag))
		}
	}

	return nil
}

func (c *Context) bindUrl(dst any) {

	v := reflect.ValueOf(dst).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if tag := field.Tag.Get("query"); tag != "" {
			v.Field(i).SetString(c.Query(tag))
		}

		if tag := field.Tag.Get("param"); tag != "" {
			v.Field(i).SetString(c.Param(tag))
		}
	}

}

func (c *Context) Param(key string) string {
	return c.params[key]
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) AddHeader(key, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) FormValues(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
}

package mapes

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	*http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r}
}

func (c *Context) Json(value interface{}) error {
	indent, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	_, err = c.W.Write([]byte(indent))
	return err
}

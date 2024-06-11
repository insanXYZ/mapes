package mapes

import (
	"fmt"
	"net/http"
)

type Handler func(ctx *Context)

type Config struct {
	Handler Handler
	Method  string
}

type Mapes struct {
	Routes map[string]Config
}

func New() *Mapes {
	return &Mapes{
		Routes: make(map[string]Config),
	}
}

func (m *Mapes) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if f, ok := m.Routes[request.URL.Path]; ok && f.Method == request.Method {
		context := NewContext(writer, request)
		f.Handler(context)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Not found"))
	}

}

func (m *Mapes) Start(port string) error {
	fmt.Println("Starting mapes on port", port)
	return http.ListenAndServe("localhost:"+port, m)
}

func (m *Mapes) add(method string, path string, handler Handler) {
	m.Routes[path] = Config{
		Handler: handler,
		Method:  method,
	}
}

func (m *Mapes) Get(route string, handler Handler) {
	m.add(http.MethodGet, route, handler)
}

func (m *Mapes) Post(route string, handler Handler) {
	m.add(http.MethodPost, route, handler)
}

func (m *Mapes) Delete(route string, handler Handler) {
	m.add(http.MethodDelete, route, handler)
}

func (m *Mapes) Put(route string, handler Handler) {
	m.add(http.MethodPut, route, handler)
}

func (m *Mapes) Patch(route string, handler Handler) {
	m.add(http.MethodPatch, route, handler)
}

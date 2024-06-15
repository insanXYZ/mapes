package mapes

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler func(ctx *Context) error

type Config struct {
	handler     Handler
	method      string
	middlewares []Handler
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

	for pattern, config := range m.Routes {
		if match, params := m.matchPattern(pattern, request.URL.Path); match && config.method == request.Method {
			context := NewContext(writer, request)
			context.params = params

			for _, middleware := range config.middlewares {
				err := middleware(context)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
					writer.Write([]byte(err.Error()))
					return
				}
			}

			config.handler(context)
			return
		}
	}

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("Not found"))

}

func (m *Mapes) Start(port string) error {
	fmt.Println("Starting mapes on port", port)
	return http.ListenAndServe("localhost:"+port, m)
}

func (m *Mapes) matchPattern(pattern, path string) (bool, map[string]string) {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			params[patternParts[i][1:]] = pathParts[i]
		} else if patternParts[i] != pathParts[i] {
			return false, nil
		}
	}

	return true, params
}

//Route

func (m *Mapes) add(method string, path string, handler Handler, middlewares []Handler) {
	m.Routes[path] = Config{
		handler:     handler,
		method:      method,
		middlewares: middlewares,
	}
}

func (m *Mapes) Get(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodGet, route, handler, middlewares)
}

func (m *Mapes) Post(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodPost, route, handler, middlewares)
}

func (m *Mapes) Delete(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodDelete, route, handler, middlewares)
}

func (m *Mapes) Put(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodPut, route, handler, middlewares)
}

func (m *Mapes) Patch(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodPatch, route, handler, middlewares)
}

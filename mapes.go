package mapes

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type (
	Handler           func(ctx *Context) error
	MiddlewareHandler func(next Handler) Handler
)

type RouteConfig struct {
	handler     Handler
	method      string
	middlewares []MiddlewareHandler
}

type RouterGroup struct {
	prefix      string
	middlewares []MiddlewareHandler
	maps        *Mapes
}

type Mapes struct {
	routes      map[string]RouteConfig
	middlewares []MiddlewareHandler
}

func New() *Mapes {
	return &Mapes{
		routes: make(map[string]RouteConfig),
	}
}

func (m *Mapes) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for pattern, config := range m.routes {
		if match, params := m.matchPattern(pattern, request.URL.Path); match && config.method == request.Method {
			context := NewContext(writer, request)
			context.params = params

			handlers := config.handler

			for _, middleware := range config.middlewares {
				handlers = middleware(handlers)
			}

			err := handlers(context)
			if err != nil {
				context.String(500, "server error")
			}

			return

		}
	}

	writer.WriteHeader(http.StatusNotFound)
	writer.Write([]byte("Not found"))
}

func (m *Mapes) Start(port string) error {
	fmt.Println("_  _ ____ ___  ____ ____\n|\\/| |__| |__] |___ [__ \n|  | |  | |    |___ ___]")
	fmt.Println("Github : https://github.com/insanXYZ/mapes")
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	err = listen.Close()
	if err != nil {
		return err
	}

	fmt.Println("Starting mapes on port", port)
	return http.ListenAndServe("localhost"+port, m)
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

func (m *Mapes) Use(middlewares ...MiddlewareHandler) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *Mapes) add(method string, path string, handler Handler, middlewares []MiddlewareHandler) {
	m.routes[path] = RouteConfig{
		handler:     handler,
		method:      method,
		middlewares: middlewares,
	}
}

func (m *Mapes) Get(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodGet, path, handler, middlewares)
}

func (m *Mapes) Post(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodPost, path, handler, middlewares)
}

func (m *Mapes) Options(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodOptions, path, handler, middlewares)
}

func (m *Mapes) Head(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodHead, path, handler, middlewares)
}

func (m *Mapes) Delete(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodDelete, path, handler, middlewares)
}

func (m *Mapes) Put(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodPut, path, handler, middlewares)
}

func (m *Mapes) Patch(path string, handler Handler, middlewares ...MiddlewareHandler) {
	m.add(http.MethodPatch, path, handler, middlewares)
}

func (m *Mapes) Static(path string, fsRoot string, middlewares ...MiddlewareHandler) {
	dir := http.Dir("./" + fsRoot)

	if string(path[len(path)-1]) == "/" {
		path = path[:len(path)-1]
	}

	fileServer := http.StripPrefix(path, http.FileServer(dir))
	staticHandler := func(ctx *Context) error {
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		return nil
	}

	if string(path[len(path)-1]) != "/" {
		path += "/"
	}

	m.Get(path+":_", staticHandler, middlewares...)
}

func (m *Mapes) Group(pattern string) *RouterGroup {
	return &RouterGroup{
		prefix:      pattern,
		middlewares: m.middlewares,
		maps:        m,
	}
}

func (rt *RouterGroup) add(method string, path string, handler Handler, middlewares []MiddlewareHandler) {
	rt.maps.routes[rt.prefix+path] = RouteConfig{
		handler:     handler,
		method:      method,
		middlewares: append(rt.middlewares, middlewares...),
	}
}

func (rt *RouterGroup) Use(middlewares ...MiddlewareHandler) {
	rt.middlewares = append(rt.middlewares, middlewares...)
}

func (rt *RouterGroup) Get(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodGet, path, handler, middlewares)
}

func (rt *RouterGroup) Options(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodOptions, path, handler, middlewares)
}

func (rt *RouterGroup) Head(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodHead, path, handler, middlewares)
}

func (rt *RouterGroup) Post(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodPost, path, handler, middlewares)
}

func (rt *RouterGroup) Delete(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodDelete, path, handler, middlewares)
}

func (rt *RouterGroup) Put(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodPut, path, handler, middlewares)
}

func (rt *RouterGroup) Patch(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodPatch, path, handler, middlewares)
}

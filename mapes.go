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

type routeConfig struct {
	handler     Handler
	method      string
	middlewares []MiddlewareHandler
}

type routerGroup struct {
	prefix      string
	middlewares []MiddlewareHandler
	maps        *Mapes
}

type Mapes struct {
	routes      map[string]routeConfig
	middlewares []MiddlewareHandler
}

func New() *Mapes {
	return &Mapes{
		routes: make(map[string]routeConfig),
	}
}

func (m *Mapes) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := NewContext(writer, request)

	for pattern, config := range m.routes {
		if match, params := m.matchPattern(pattern, request.URL.Path); match && config.method == request.Method {
			context.params = params
			handlers := config.handler

			for _, middleware := range config.middlewares {
				handlers = middleware(handlers)
			}

			if err := handlers(context); err != nil {
				context.Json(500, map[string]string{
					"message": "server error",
				})
			}
			return
		}
	}

	context.Json(404, map[string]string{
		"message": "not found",
	})
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
	m.routes[path] = routeConfig{
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

func (m *Mapes) Group(pattern string) *routerGroup {
	return &routerGroup{
		prefix:      pattern,
		middlewares: m.middlewares,
		maps:        m,
	}
}

func (rt *routerGroup) add(method string, path string, handler Handler, middlewares []MiddlewareHandler) {
	rt.maps.routes[rt.prefix+path] = routeConfig{
		handler:     handler,
		method:      method,
		middlewares: append(rt.middlewares, middlewares...),
	}
}

func (rt *routerGroup) Use(middlewares ...MiddlewareHandler) {
	rt.middlewares = append(rt.middlewares, middlewares...)
}

func (rt *routerGroup) Get(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodGet, path, handler, middlewares)
}

func (rt *routerGroup) Options(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodOptions, path, handler, middlewares)
}

func (rt *routerGroup) Head(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodHead, path, handler, middlewares)
}

func (rt *routerGroup) Post(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodPost, path, handler, middlewares)
}

func (rt *routerGroup) Delete(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodDelete, path, handler, middlewares)
}

func (rt *routerGroup) Put(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodPut, path, handler, middlewares)
}

func (rt *routerGroup) Patch(path string, handler Handler, middlewares ...MiddlewareHandler) {
	rt.add(http.MethodPatch, path, handler, middlewares)
}

package mapes

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type Handler func(ctx *Context) error

type RouteConfig struct {
	handler     Handler
	method      string
	middlewares []Handler
}

type RouterGroup struct {
	prefix      string
	middlewares []Handler
	maps        *Mapes
}

type Mapes struct {
	routes      map[string]RouteConfig
	middlewares []Handler
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

//Route

func (m *Mapes) Use(middlewares ...Handler) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *Mapes) add(method string, path string, handler Handler, middlewares []Handler) {
	m.routes[path] = RouteConfig{
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

func (m *Mapes) Options(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodOptions, route, handler, middlewares)
}

func (m *Mapes) Head(route string, handler Handler, middlewares ...Handler) {
	m.add(http.MethodHead, route, handler, middlewares)
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

//Router Group

func (m *Mapes) Group(pattern string) *RouterGroup {
	return &RouterGroup{
		prefix:      pattern,
		middlewares: m.middlewares,
		maps:        m,
	}
}

func (rt *RouterGroup) add(method string, path string, handler Handler, middlewares []Handler) {
	rt.maps.routes[rt.prefix+path] = RouteConfig{
		handler:     handler,
		method:      method,
		middlewares: append(rt.middlewares, middlewares...),
	}
}

func (rt *RouterGroup) Use(middlewares ...Handler) {
	rt.middlewares = append(rt.middlewares, middlewares...)
}

func (rt *RouterGroup) Get(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodGet, route, handler, middlewares)
}

func (rt *RouterGroup) Options(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodOptions, route, handler, middlewares)
}

func (rt *RouterGroup) Head(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodHead, route, handler, middlewares)
}

func (rt *RouterGroup) Post(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodPost, route, handler, middlewares)
}

func (rt *RouterGroup) Delete(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodDelete, route, handler, middlewares)
}

func (rt *RouterGroup) Put(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodPut, route, handler, middlewares)
}

func (rt *RouterGroup) Patch(route string, handler Handler, middlewares ...Handler) {
	rt.add(http.MethodPatch, route, handler, middlewares)
}

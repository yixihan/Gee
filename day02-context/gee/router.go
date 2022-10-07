package gee

import (
	"log"
	"net/http"
)

type router struct {
	handler map[string]HandlerFunc
}

// NewRouter router 构造函数
func NewRouter() *router {
	return &router{
		handler: make(map[string]HandlerFunc),
	}
}

// addRoute 添加一个新的路由到 router 中
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handler[key] = handler
}

// handle 路由
func (r *router) handle(c *Context) {
	log.Printf("Route %4s - %s", c.Method, c.Path)
	key := c.Method + "-" + c.Path
	if handler, ok := r.handler[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found", c.Path)
	}
}

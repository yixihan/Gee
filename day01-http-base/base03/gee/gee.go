package gee

import (
	"fmt"
	"net/http"
)

//HandlerFunc 定义了 gee 使用的请求处理程序
type HandlerFunc func(http.ResponseWriter, *http.Request)

//Engine 实现了 ServeHandler 接口
type Engine struct {
	router map[string]HandlerFunc
}

//New Engine 的构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

//addRoute Engine 添加路由的函数, 私有
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	e.router[key] = handler
}

//GET 定义添加 GET 请求的方法
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

//POST 定义添加 POST 请求的方法
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

//Run 定义启动 http 服务器的方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

//ServeHTTP 自定义实现的的 ServeHTTP 方法
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if hander, ok := e.router[key]; ok {
		hander(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found : %s\n", r.URL.Path)
	}
}

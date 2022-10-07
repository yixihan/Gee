package gee

import "net/http"

// HandlerFunc 定义了 gee 使用的请求处理程序, 参数为自定义的 Context
type HandlerFunc func(*Context)

// Engine 实现 ServeHTTP 接口
type Engine struct {
	router *router
}

// New Engine 构造函数
func New() *Engine {
	return &Engine{
		router: NewRouter(),
	}
}

// addRoute 添加一个新的路由
func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
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

//ServeHTTP 自定义实现的的 ServeHTTP 方法, 具体处理逻辑放在 router.go
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(r, w)
	e.router.handle(c)
}

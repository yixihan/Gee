package gee

import (
	"log"
	"net/http"
)

// HandlerFunc 定义了 gee 使用的请求处理程序, 参数为自定义的 Context
type HandlerFunc func(*Context)

type (
	// Engine 实现 ServeHTTP 接口
	Engine struct {
		// RouterGroup 路由组指针
		*RouterGroup
		// router 路由
		router *router
		// groups 路由组数组
		groups []*RouterGroup
	}

	// RouterGroup 路由组
	RouterGroup struct {
		// prefix 前缀
		prefix string
		// middlewares 中间件
		middlewares []HandlerFunc
		// parent 父路由组
		parent *RouterGroup
		// engine engine
		engine *Engine
	}
)

// New Engine 构造函数
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group RouterGroup 构造函数
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute 添加一个新的路由
func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s\n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//GET 定义添加 GET 请求的方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

//POST 定义添加 POST 请求的方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//Run 定义启动 http 服务器的方法
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

//ServeHTTP 自定义实现的的 ServeHTTP 方法, 具体处理逻辑放在 router.go
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(r, w)
	e.router.handle(c)
	log.Print(c.Params)
}

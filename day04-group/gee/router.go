package gee

import (
	"log"
	"net/http"
	"strings"
)

// router router 定义
// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// NewRouter router 构造函数
func NewRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 路由解析
// * 为全匹配模式
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0, len(vs))
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

// addRoute 添加一个路由进 router
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	// 获取 parts
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	// 如果 roots 里面没有当前请求方法的数, 则新建一颗树
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	// roots 中插入新节点, handler 中放入新 handler
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute 获取路由 获取路由对应树节点和 params
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	// 获取 path 对应的 parts
	searchParts := parsePattern(path)
	params := make(map[string]string)

	// 尝试从 roots 里面获取对应方法的 trie 树
	root, ok := r.roots[method]

	// 如果没有, 返回 nil
	if !ok {
		return nil, nil
	}

	// 如果有的话, 从树中搜索节点
	if n := root.search(searchParts, 0); n != nil {
		// 如果搜索到节点, 则将节点中的 pattern 解析出来
		parts := parsePattern(n.pattern)

		for index, part := range parts {
			// 如果是全匹配模式 (*)
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// 如果是动态路由
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

// handle 路由匹配
func (r *router) handle(c *Context) {
	log.Printf("Route %4s - %s\n", c.Method, c.Path)
	if n, params := r.getRoute(c.Method, c.Path); n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found : %s\n", c.Path)
	}
}

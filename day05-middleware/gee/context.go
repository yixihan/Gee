package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H 给 map[string]interface{} 起了一个别名 gee.H, 构建 JSON 数据时, 显得更简洁
type H map[string]interface{}

// Context context 定义
type Context struct {
	// Req request
	Req *http.Request

	// Writer responseWriter
	Writer http.ResponseWriter

	// Path 路由路径
	Path string

	// Method 请求方法
	Method string

	// Params 参数
	Params map[string]string

	// StatusCode 响应状态码
	StatusCode int

	// handlers 中间件
	handlers []HandlerFunc

	// index 记录当前执行到第几个中间件
	index int
}

// NewContext Context 构造函数
func NewContext(req *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Req:    req,
		Writer: w,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next 执行下一个中间件
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)

	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// PostForm 根据 key 返回 form 中对应的 value
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 根据 key 返回 query 中对应的 value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Param 根据 key 返回 params 中对应的 value
func (c *Context) Param(key string) string {
	if value, ok := c.Params[key]; !ok {
		return ""
	} else {
		return value
	}
}

// Status 设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(c.StatusCode)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

// String 响应体以 string 格式返回
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 响应体以 JSON 格式返回
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data 响应体为 data 数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// HTML 响应体以 HTML 格式返回
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	_, _ = c.Writer.Write([]byte(html))
}

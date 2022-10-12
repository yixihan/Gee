package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// stace 打印错误信息
func trace(message string) string {
	pcs := make([]uintptr, 32)
	// 跳过前三个 caller
	n := runtime.Callers(3, pcs[:])

	str := new(strings.Builder)
	str.WriteString(message + "\nTraceBack:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}

	return str.String()
}

// Recovery 错误恢复中间件
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
	}
}

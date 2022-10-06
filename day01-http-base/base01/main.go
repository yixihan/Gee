package main

import (
	"fmt"
	"log"
	"net/http"
)

//main
func main() {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	/*
		http.ListenAndServe(addr, handler)
		作用 : 启动 web 服务
		参数解释 :
			addr : 地址, :9999 表示在 9999 端口监听
			handler : 代表处理所有的HTTP请求的实例, nil 代表使用标准库中的实例处理
		第二个参数, 则是我们基于 net / http 标准库实现 web 框架的入口
	*/
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	for k, v := range r.Header {
		_, _ = fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

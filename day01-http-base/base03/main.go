package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	r := gee.New()

	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "URL,Path = %q\n", r.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		for key, value := range r.Header {
			_, _ = fmt.Fprintf(w, "Header[%q]: %q\n", key, value)
		}
	})

	_ = r.Run(":9999")
}

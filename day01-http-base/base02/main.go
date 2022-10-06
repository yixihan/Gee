package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		_, _ = fmt.Fprintf(w, "RL.Path = %q\n", r.URL.Path)
	case "/hello":
		for key, val := range r.Header {
			_, _ = fmt.Fprintf(w, "Header[%q] = %q\n", key, val)
		}
	default:
		_, _ = fmt.Fprintf(w, "404 Not Found : %s\n", r.URL.Path)
	}
}

func main() {
	engine := &Engine{}
	log.Fatal(http.ListenAndServe(":9999", engine))
}

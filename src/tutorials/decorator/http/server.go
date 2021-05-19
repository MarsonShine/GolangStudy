package http

import (
	"fmt"
	"log"
	"net/http"
)

func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Println("---->WithServerHeader()")
		rw.Header().Set("Server", "HelloServer v0.0.1")
		h(rw, r) // next
	}
}
func Hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}

// Pipeline

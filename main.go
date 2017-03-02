package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type fileHandler struct {
	handle http.Handler
}

func (h fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	begin := time.Now()
	log.Printf("[Remote=%s] [Path=%s] [RequestStart]", r.RemoteAddr, r.URL.Path)
	w.Header().Set("Cache-control", "no-store")
	h.handle.ServeHTTP(w, r)
	log.Printf("[Remote=%s] [Path=%s] [RequestEnd] [ResponseTime=%s]",
		r.RemoteAddr, r.URL.Path, time.Since(begin))
}

func main() {
	port := flag.String("port", "8085", "Port to listen")
	path := flag.String("path", ".", "path to serve")
	flag.Parse()

	log.Println("Listenting on port", *port, "serving", *path)

	fs := http.FileServer(http.Dir(*path))
	http.Handle("/", fileHandler{http.StripPrefix("/", fs)})
	err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

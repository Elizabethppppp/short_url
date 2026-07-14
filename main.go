package main

import (
	"log"
	"sync"

	server "github.com/Elizabethppppp/tcp_server"
)

func main() {
	store := NewURLstore()

	mux := server.NewMux()
	mux.Handle("/hello", func(w server.ResponseWriter, r *server.Request) {
		w.WriteHeader(200)
		w.SetHeader("Content-Type", "text/plain")
		w.Write([]byte("Hello from imported lib!"))
	})

	if err := server.Listen(":8090", mux); err != nil {
		log.Fatal(err)
	}
}

type URLstore struct {
	mu    sync.RWMutex
	links map[string]string
	count map[string]int
}

func NewURLstore() *URLstore {
	return &URLstore{
		links: make(map[string]string),
		count: make(map[string]int),
	}
}

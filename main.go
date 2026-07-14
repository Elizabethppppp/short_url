package main

import (
	"log"

	server "github.com/Elizabethppppp/tcp_server"
)

func main() {
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

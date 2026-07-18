package main

import (
	"log"

	server "github.com/Elizabethppppp/tcp_server"
)

func main() {
	store := NewURLstore()

	mux := server.NewMux()
	mux.Handle("/short", store.CreateShortURL)
	mux.Handle("/{short}", store.RedirectHandler)
	mux.Handle("/count/{short}", store.CountShortURL)

	if err := server.Listen(":8090", mux); err != nil {
		log.Fatal(err)
	}
}

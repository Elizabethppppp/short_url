package main

import (
	"fmt"
	"strings"

	server "github.com/Elizabethppppp/tcp_server"
)

type URLstore struct {
	links   map[string]string
	count   map[string]int
	counter uint64
}

func NewURLstore() *URLstore {
	return &URLstore{
		links:   make(map[string]string),
		count:   make(map[string]int),
		counter: 100000000000,
	}
}

// post method
func (u *URLstore) CreateShortURL(w server.ResponseWriter, r *server.Request) {
	if r.Method != "POST" {
		w.WriteHeader(server.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	originalURL := strings.TrimSpace(string(r.Body))

	if originalURL == "" {
		w.WriteHeader(server.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	var shortURL string
	var inMap bool

	for short, original := range u.links {
		if original == originalURL {
			shortURL = short
			inMap = true
			break
		}
	}

	if inMap {
		response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURL)
		w.WriteHeader(server.StatusOK)
		w.Write([]byte(response))
		return
	}

	shortURL = u.generateShortURL()

	for {
		if _, in := u.links[shortURL]; !in {
			break
		}
		shortURL = u.generateShortURL()
	}

	u.links[shortURL] = originalURL
	u.count[shortURL] = 0

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURL)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

// get method
func (u *URLstore) RedirectHandler(w server.ResponseWriter, r *server.Request) {
	if r.Method != "GET" {
		w.WriteHeader(server.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	shortURL := r.Param("short")

	originalURL, inMap := u.links[shortURL]

	if !inMap {
		w.WriteHeader(server.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	u.count[shortURL]++

	w.SetHeader("Location", originalURL)
	w.WriteHeader(server.StatusMoving)
	w.Write([]byte("Redirecting to " + originalURL))
}

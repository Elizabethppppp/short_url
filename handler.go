package main

import (
	"fmt"
	"net/url"
	"strings"

	server "github.com/Elizabethppppp/tcp_server"
)

type URLstore struct {
	links map[string]string
	count map[string]int
}

func NewURLstore() *URLstore {
	return &URLstore{
		links: make(map[string]string),
		count: make(map[string]int),
	}
}

// post method
func (u *URLstore) CreateShortURL(w server.ResponseWriter, r *server.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	originalURL := strings.TrimSpace(string(r.Body))

	if originalURL == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}

	parsedURL, err := url.ParseRequestURI(originalURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		w.WriteHeader(400)
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
		response := fmt.Sprintf(`{"shortURL":"http://localhost:8080/%s"}`, shortURL)
		w.WriteHeader(200)
		w.Write([]byte(response))
		return
	}

	shortURL = generateShortURL(originalURL)

	for {
		if _, in := u.links[shortURL]; !in {
			break
		}
		shortURL = generateShortURL(originalURL + shortURL)
	}

	u.links[shortURL] = originalURL
	u.count[shortURL] = 0

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8080/%s"}`, shortURL)
	w.WriteHeader(200)
	w.Write([]byte(response))
}

func generateShortURL(originalURL string) string {
	return originalURL
}

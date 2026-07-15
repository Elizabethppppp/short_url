package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	server "github.com/Elizabethppppp/tcp_server"
	"github.com/yihleego/base62"
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

	shortURL = u.generateShortURL()

	for {
		if _, in := u.links[shortURL]; !in {
			break
		}
		shortURL = u.generateShortURL()
	}

	u.links[shortURL] = originalURL
	u.count[shortURL] = 0

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8080/%s"}`, shortURL)
	w.WriteHeader(200)
	w.Write([]byte(response))
}

func (u *URLstore) generateShortURL() string {
	u.counter++

	str := strconv.FormatUint(u.counter, 10)
	encoded := base62.StdEncoding.Encode([]byte(str))

	result := string(encoded)
	if len(result) < 7 {
		result = "0" + result
	}

	return result
}

// get method
func (u *URLstore) RedirectHandler(w server.ResponseWriter, r *server.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	path := r.Path

	shortURL := strings.TrimPrefix(path, "/")

	if shortURL == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}

	originalURL, inMap := u.links[shortURL]

	if !inMap {
		w.WriteHeader(404)
		w.Write([]byte("Not Found"))
		return
	}

	u.count[shortURL]++

	w.SetHeader("Location", originalURL)
	w.WriteHeader(302)
	w.Write([]byte("Redirecting to " + originalURL))
}

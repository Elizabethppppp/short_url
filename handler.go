package main

import (
	"context"
	"fmt"
	"strings"

	server "github.com/Elizabethppppp/tcp_server"
	"github.com/jackc/pgx/v5"
)

type URLstore struct {
	db *pgx.Conn
}

func NewURLstore(db *pgx.Conn) *URLstore {
	return &URLstore{
		db: db,
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

	ctx := context.Background()

	var shortURLdb string
	err := u.db.QueryRow(ctx, "SELECT shortURL FROM url_schema.url WHERE originalURL = $1", originalURL).Scan(&shortURLdb)
	if err == nil {
		response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURLdb)
		w.WriteHeader(server.StatusOK)
		w.Write([]byte(response))
		return
	}

	if err != pgx.ErrNoRows {
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Database Error"))
		return
	}

	shortURL, counter, err := u.generateShortURL(ctx)
	if err != nil {
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Counter Error"))
		return
	}

	_, err = u.db.Exec(ctx, "INSERT INTO url_schema.url (originalURL, shortURL, count, last_counter) VALUES ($1, $2, 0, $3)",
		originalURL, shortURL, counter)

	if err != nil {
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Insert error"))
		return
	}

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURL)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

// get method
/*func (u *URLstore) RedirectHandler(w server.ResponseWriter, r *server.Request) {
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

// get method for count
func (u *URLstore) CountShortURL(w server.ResponseWriter, r *server.Request) {
	if r.Method != "GET" {
		w.WriteHeader(server.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	shortURL := r.Param("short")

	_, inMap := u.links[shortURL]
	if !inMap {
		w.WriteHeader(server.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	c := u.count[shortURL]

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s", "count":%d}`, shortURL, c)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}*/

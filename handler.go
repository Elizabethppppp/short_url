package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"test/logger"

	server "github.com/Elizabethppppp/tcp_server"
)

type URLStore struct {
	db *sql.DB
	mu sync.Mutex
}

func NewURLstore(db *sql.DB) *URLStore {
	return &URLStore{
		db: db,
	}
}

// post method
func (u *URLStore) CreateShortURL(w server.ResponseWriter, r *server.Request) {

	originalURL, err1 := server.ParseURL(string(r.Body))
	if err1 != nil {
		ResponseJSON(w, 400, err1)
		return
	}

	if originalURL == nil {
		w.WriteHeader(server.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	ctx := context.Background()

	var shortURLdb string
	err := u.db.QueryRowContext(ctx, "SELECT shortURL FROM url WHERE originalURL = $1", originalURL.String()).Scan(&shortURLdb)
	if err == nil {
		response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURLdb)
		w.WriteHeader(server.StatusOK)
		w.Write([]byte(response))
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		ResponseJSON(w, 500, err)
		return
	}

	shortURL, counter, err := u.generateShortURL(ctx)
	if err != nil {
		ResponseJSON(w, 500, err)
		return
	}

	_, err = u.db.ExecContext(ctx, "INSERT INTO url (originalURL, shortURL, count, last_counter) VALUES ($1, $2, 0, $3)",
		originalURL.String(), shortURL, counter)

	if err != nil {
		ResponseJSON(w, 500, err)
		return
	}

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURL)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

// get method
func (u *URLStore) RedirectHandler(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	var originalURL string
	err := u.db.QueryRowContext(ctx, "SELECT originalURL FROM url WHERE shortURL = $1", shortURL).Scan(&originalURL)
	if errors.Is(err, sql.ErrNoRows) {
		ResponseJSON(w, 404, err)
		return
	}
	if err != nil {
		ResponseJSON(w, 500, err)
		return
	}

	_, err1 := u.db.ExecContext(ctx, "UPDATE url SET count = count + 1 WHERE shortURL = $1", shortURL)
	if err1 != nil {
		ResponseJSON(w, 500, err1)
		return
	}

	w.SetHeader("Location", originalURL)
	w.WriteHeader(server.StatusMoving)
	w.Write([]byte("Redirecting to " + originalURL))
}

// get method for count
func (u *URLStore) CountShortURL(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	var count int
	err := u.db.QueryRowContext(ctx, "SELECT count FROM url WHERE shortURL = $1", shortURL).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		ResponseJSON(w, 404, err)
		return
	}
	if err != nil {
		logger.Error("Count Error", "shortURL", shortURL, "error", err)
		ResponseJSON(w, 500, err)
		return
	}

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s", "count":%d}`, shortURL, count)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"test/logger"

	server "github.com/Elizabethppppp/tcp_server"
	"github.com/jackc/pgx/v5"
)

type URLstore struct {
	db *sql.DB
}

func NewURLstore(db *sql.DB) *URLstore {
	return &URLstore{
		db: db,
	}
}

// post method
func (u *URLstore) CreateShortURL(w server.ResponseWriter, r *server.Request) {

	logger.Debug("CreateShortURL", "method", r.Method, "path", r.Path)

	if r.Method != "POST" {
		logger.Warn("Method Not Allowed", "method", r.Method, "expected", "POST")
		w.WriteHeader(server.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	originalURL := strings.TrimSpace(string(r.Body))

	if originalURL == "" {
		logger.Warn("Empty URL")
		w.WriteHeader(server.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	ctx := context.Background()

	var shortURLdb string
	err := u.db.QueryRowContext(ctx, "SELECT shortURL FROM url_schema.url WHERE originalURL = $1", originalURL).Scan(&shortURLdb)
	if err == nil {
		logger.Info("Return shortURL", "originalURL", originalURL, "shortURL", shortURLdb)
		response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURLdb)
		w.WriteHeader(server.StatusOK)
		w.Write([]byte(response))
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		logger.Error("Database Error", "originalURL", originalURL, "error", err)
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Database Error"))
		return
	}

	shortURL, counter, err := u.generateShortURL(ctx)
	if err != nil {
		logger.Error("Generate shortUrl Error", "originalURL", originalURL, "error", err)
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Counter Error"))
		return
	}

	_, err = u.db.ExecContext(ctx, "INSERT INTO url_schema.url (originalURL, shortURL, count, last_counter) VALUES ($1, $2, 0, $3)",
		originalURL, shortURL, counter)

	if err != nil {
		logger.Error("Insert Error","error", err, "originalURL", originalURL,"shortURL", shortURL )
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Insert error"))
		return
	}

	logger.Info("Insert Success", "originalURL", originalURL, "shortURL", shortURL ,"counter", counter)

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURL)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

// get method
func (u *URLstore) RedirectHandler(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	logger.Debug("RedirectHandler", "method", r.Method, "path", r.Path, "ShortURL", shortURL)

	if r.Method != "GET" {
		logger.Warn("Method Not Allowed", "method", r.Method, "expected", "GET")
		w.WriteHeader(server.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	var originalURL string
	err := u.db.QueryRowContext(ctx, "SELECT originalURL FROM url_schema.url WHERE shortURL = $1", shortURL).Scan(&originalURL)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Warn("ShortUrl not found",  "ShortURL", shortURL)
		w.WriteHeader(server.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	if err != nil {
		logger.Error("Insert Error", "shortURL", shortURL, "error", err)
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Insert Error"))
		return
	}

	_, err1 := u.db.ExecContext(ctx, "UPDATE url_schema.url SET count = count + 1 WHERE shortURL = $1", shortURL)
	if err1 != nil {
		logger.Error("Update Error", "shortURL", shortURL, "error", err1)
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Update Error"))
		return
	}

	logger.Info("Success", "shortURL", shortURL, "originalURL", originalURL)

	w.SetHeader("Location", originalURL)
	w.WriteHeader(server.StatusMoving)
	w.Write([]byte("Redirecting to " + originalURL))
}

// get method for count
func (u *URLstore) CountShortURL(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	logger.Debug("CountShortURL", "method", r.Method, "path", r.Path, "ShortURL", shortURL)

	if r.Method != "GET" {
		logger.Warn("Method Not Allowed", "method", r.Method, "expected", "GET")
		w.WriteHeader(server.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	var count int
	err := u.db.QueryRowContext(ctx, "SELECT count FROM url_schema.url WHERE shortURL = $1", shortURL).Scan(&count)
	if errors.Is(err, pgx.ErrNoRows) {
		logger.Warn("ShortUrl not found",  "ShortURL", shortURL)
		w.WriteHeader(server.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
	if err != nil {
		logger.Error("Count Error", "shortURL", shortURL, "error", err)
		w.WriteHeader(server.StatusInternalServerError)
		w.Write([]byte("Counter Error"))
		return
	}

	logger.Debug("Count Success", "shortURL", shortURL, "count", count)

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s", "count":%d}`, shortURL, count)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

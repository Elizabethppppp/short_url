package handler

import (
	"context"
	"errors"
	"fmt"
	"test/error_response"
	"test/logger"
	"test/pgService"

	server "github.com/Elizabethppppp/tcp_server"
)

type URLStore struct {
	pg *pgService.PgService
}

func NewURLstore(pg *pgService.PgService) *URLStore {
	return &URLStore{
		pg: pg,
	}
}

// post method
func (u *URLStore) CreateShortURL(w server.ResponseWriter, r *server.Request) {

	originalURL, err1 := server.ParseURL(string(r.Body))
	if err1 != nil {
		error_response.ResponseJSON(w, 400, err1)
		return
	}

	if originalURL == nil {
		w.WriteHeader(server.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	}

	ctx := context.Background()

	shortURLdb, err := u.pg.GetShortURL(ctx, originalURL.Raw)
	if err == nil {
		response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s"}`, shortURLdb)
		w.WriteHeader(server.StatusOK)
		w.Write([]byte(response))
		return
	}

	if errors.Is(err, pgService.ErrInternal) {
		error_response.ResponseJSON(w, 500, err)
		return
	}

	shortURL, counter, err := u.generateShortURL(ctx)
	if err != nil {
		error_response.ResponseJSON(w, 500, err)
		return
	}

	if err := u.pg.Insert(ctx, originalURL.Raw, shortURL, counter); err != nil {
		error_response.ResponseJSON(w, 500, err)
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

	originalURL, err := u.pg.RedirectShortURL(ctx, shortURL)
	if errors.Is(err, pgService.ErrNotFound) {
		error_response.ResponseJSON(w, 404, err)
		return
	}
	if err != nil {
		error_response.ResponseJSON(w, 500, err)
		return
	}

	if err1 := u.pg.UpdateCounter(ctx, shortURL); err1 != nil {
		error_response.ResponseJSON(w, 500, err1)
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

	count, err := u.pg.GetCount(ctx, shortURL)
	if errors.Is(err, pgService.ErrNotFound) {
		error_response.ResponseJSON(w, 404, err)
		return
	}
	if err != nil {
		logger.Error("Count Error", "shortURL", shortURL, "error", err)
		error_response.ResponseJSON(w, 500, err)
		return
	}

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s", "count":%d}`, shortURL, count)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

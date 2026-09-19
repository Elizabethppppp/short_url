package handler

import (
	"context"
	"errors"
	"fmt"
	"test/error_response"
	"test/logger"
	"test/reduceService"

	server "github.com/Elizabethppppp/tcp_server"
)

type URLStore struct {
	s *reduceService.ReduceService
}

func NewURLstore(s *reduceService.ReduceService) *URLStore {
	return &URLStore{
		s: s,
	}
}

// post method
func (u *URLStore) CreateShortURL(w server.ResponseWriter, r *server.Request) {

	ctx := context.Background()
	shortURL, err := u.s.Create(ctx, string(r.Body))
	if errors.Is(err, reduceService.ErrBadRequest) {
		error_response.ResponseJSON(w, server.StatusBadRequest, err)
		return
	}
	if err != nil {
		error_response.ResponseJSON(w, server.StatusInternalServerError, err)
		return
	}

	response := fmt.Sprintf(`{"shortURL":"%q"}`, shortURL)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

// get method
func (u *URLStore) RedirectHandler(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	originalURL, err := u.s.Redirect(ctx, shortURL)
	if errors.Is(err, reduceService.ErrNotFound) {
		error_response.ResponseJSON(w, server.StatusNotFound, err)
		return
	}
	if err != nil {
		error_response.ResponseJSON(w, server.StatusInternalServerError, err)
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

	_, count, err := u.s.Count(ctx, shortURL)
	if errors.Is(err, reduceService.ErrNotFound) {
		error_response.ResponseJSON(w, server.StatusNotFound, err)
		return
	}
	if err != nil {
		logger.Error("Count Error", "shortURL", shortURL, "error", err)
		error_response.ResponseJSON(w, server.StatusInternalServerError, err)
		return
	}

	response := fmt.Sprintf(`{"shortURL":"http://localhost:8090/%s", "count":%d}`, shortURL, count)
	w.WriteHeader(server.StatusOK)
	w.Write([]byte(response))
}

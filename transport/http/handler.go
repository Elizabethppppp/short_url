package transport

import (
	"context"
	"errors"
	"fmt"
	"test/core/url_short"
	"test/error_response"
	"test/logger"

	server "github.com/Elizabethppppp/tcp_server"
)

// post method
func (t *Transport) CreateShortURL(w server.ResponseWriter, r *server.Request) {

	ctx := context.Background()
	shortURL, err := t.tr.Create(ctx, string(r.Body))
	if errors.Is(err, url_short.ErrBadRequest) {
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
func (t *Transport) RedirectHandler(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	originalURL, err := t.tr.Redirect(ctx, shortURL)
	if errors.Is(err, url_short.ErrNotFound) {
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
func (t *Transport) CountShortURL(w server.ResponseWriter, r *server.Request) {
	ctx := context.Background()

	shortURL := r.Param("short")

	_, count, err := t.tr.Count(ctx, shortURL)
	if errors.Is(err, url_short.ErrNotFound) {
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

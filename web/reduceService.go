package web

import (
	"context"
	"errors"
	"test/pgService"

	server "github.com/Elizabethppppp/tcp_server"
)

var (
	ErrBadRequest = errors.New("Bad Request")
	ErrNotFound   = errors.New("Not Found")
	ErrInternal   = errors.New("Internal Server Error")
)

type ReduceService struct {
	pg *pgService.PgService
}

func NewReduceService(pg *pgService.PgService) *ReduceService {
	return &ReduceService{pg: pg}
}

func (service *ReduceService) Create(ctx context.Context, url string) (string, error) {

	parsed, err := server.ParseURL(url)
	if err != nil {
		return "", ErrBadRequest
	}

	if found, err := service.pg.GetShortURL(ctx, parsed.Raw); err == nil {
		return "http://localhost:8090/" + found, nil
	} else if !errors.Is(err, pgService.ErrNotFound) {
		return "", ErrInternal
	}

	code, counter, err := service.generateShortURL(ctx)
	if err != nil {
		return "", ErrInternal
	}

	if err := service.pg.Insert(ctx, parsed.Raw, code, counter); err != nil {
		return "", ErrInternal
	}
	return "http://localhost:8090/" + code, nil
}

func (service *ReduceService) Redirect(ctx context.Context, shortURL string) (string, error) {

	originalURL, err := service.pg.RedirectShortURL(ctx, shortURL)
	if errors.Is(err, pgService.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", ErrInternal
	}

	if err := service.pg.UpdateCounter(ctx, shortURL); err != nil {
		return "", ErrInternal
	}

	return originalURL, nil
}

func (service *ReduceService) Count(ctx context.Context, shortURL string) (string, int, error) {
	count, err := service.pg.GetCount(ctx, shortURL)
	if errors.Is(err, pgService.ErrNotFound) {
		return "", 0, ErrNotFound
	}
	if err != nil {
		return "", 0, ErrInternal
	}

	return "http://localhost:8090/" + shortURL, count, nil
}

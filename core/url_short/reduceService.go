package url_short

import (
	"context"
	"errors"
	"test/serviceErrors"

	server "github.com/Elizabethppppp/tcp_server"
)

type URLRepository interface {
	GetShortURL(ctx context.Context, originalURL string) (string, error)
	Insert(ctx context.Context, originalURL, shortURL string, lastCounter uint64) error
	NextCount(ctx context.Context) (uint64, error)
	RedirectShortURL(ctx context.Context, shortURL string) (string, error)
	UpdateCounter(ctx context.Context, shortURL string) error
	GetCount(ctx context.Context, shortURL string) (int, error)
}

type ReduceService struct {
	repo URLRepository
}

func NewReduceService(repo URLRepository) *ReduceService {
	return &ReduceService{repo: repo}
}

func (service *ReduceService) Create(ctx context.Context, url string) (string, error) {

	parsed, err := server.ParseURL(url)
	if err != nil {
		return "", serviceErrors.ErrBadRequest
	}

	if found, err := service.repo.GetShortURL(ctx, parsed.Raw); err == nil {
		return "http://localhost:8090/" + found, nil
	} else if !errors.Is(err, serviceErrors.ErrNotFound) {
		return "", serviceErrors.ErrInternal
	}

	code, counter, err := service.generateShortURL(ctx)
	if err != nil {
		return "", serviceErrors.ErrInternal
	}

	if err := service.repo.Insert(ctx, parsed.Raw, code, counter); err != nil {
		return "", serviceErrors.ErrInternal
	}
	return "http://localhost:8090/" + code, nil
}

func (service *ReduceService) Redirect(ctx context.Context, shortURL string) (string, error) {

	originalURL, err := service.repo.RedirectShortURL(ctx, shortURL)
	if errors.Is(err, serviceErrors.ErrNotFound) {
		return "", serviceErrors.ErrNotFound
	}
	if err != nil {
		return "", serviceErrors.ErrInternal
	}

	if err := service.repo.UpdateCounter(ctx, shortURL); err != nil {
		return "", serviceErrors.ErrInternal
	}

	return originalURL, nil
}

func (service *ReduceService) Count(ctx context.Context, shortURL string) (string, int, error) {
	count, err := service.repo.GetCount(ctx, shortURL)
	if errors.Is(err, serviceErrors.ErrNotFound) {
		return "", 0, serviceErrors.ErrNotFound
	}
	if err != nil {
		return "", 0, serviceErrors.ErrInternal
	}

	return "http://localhost:8090/" + shortURL, count, nil
}

package pgService

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal server error")
)

type PgService struct {
	db *sql.DB
}

func NewPgService(db *sql.DB) *PgService {
	return &PgService{
		db: db,
	}
}

func (pg *PgService) GetShortURL(ctx context.Context, originalURL string) (string, error) {

	var shortURL string
	err := pg.db.QueryRowContext(ctx, "SELECT shortURL FROM url WHERE originalURL = $1", originalURL).Scan(&shortURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (pg *PgService) Insert(ctx context.Context, originalURL, shortURL string, last_counter uint64) error {
	_, err := pg.db.ExecContext(ctx, "INSERT INTO url (originalURL, shortURL, count, last_counter) VALUES ($1, $2, 0, $3)",
		originalURL, shortURL, last_counter)

	if err != nil {
		return ErrInternal
	}

	return nil
}

func (pg *PgService) NextCount(ctx context.Context) (uint64, error) {
	var count uint64
	err := pg.db.QueryRowContext(ctx, "SELECT nextval('url_seq')").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (pg *PgService) RedirectShortURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := pg.db.QueryRowContext(ctx, "SELECT originalURL FROM url WHERE shortURL = $1", shortURL).Scan(&originalURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", ErrInternal
	}
	return originalURL, nil
}

func (pg *PgService) UpdateCounter(ctx context.Context, shortURL string) error {
	_, err := pg.db.ExecContext(ctx, "UPDATE url SET count = count + 1 WHERE shortURL = $1", shortURL)
	if err != nil {
		return ErrInternal
	}
	return nil
}

func (pg *PgService) GetCount(ctx context.Context, shortURL string) (int, error) {
	var count int
	err := pg.db.QueryRowContext(ctx, "SELECT count FROM url WHERE shortURL = $1", shortURL).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, ErrInternal
	}
	return count, nil
}

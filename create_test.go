package short_url

import (
	"context"
	"errors"
	"test/core/urlShort"
	"test/serviceErrors"
	"testing"
)

type URLRepositoryTest struct {
	shortByOriginal map[string]string
	failGet         bool
	failInsert      bool
	nextVal         uint64
	insertVal       map[string]string
}

func NewURLRepositoryTest() *URLRepositoryTest {
	return &URLRepositoryTest{
		shortByOriginal: map[string]string{},
		insertVal:       map[string]string{},
		nextVal:         100000000001,
	}
}

func (u *URLRepositoryTest) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	if u.failGet {
		return "", serviceErrors.ErrInternal
	}
	if short, ok := u.shortByOriginal[originalURL]; ok {
		return short, nil
	}
	return "", serviceErrors.ErrNotFound
}

func (u *URLRepositoryTest) Insert(ctx context.Context, originalURL, shotUrl string, last_counter uint64) error {
	if u.failInsert {
		return serviceErrors.ErrInternal
	}
	u.insertVal[originalURL] = shotUrl
	u.shortByOriginal[originalURL] = shotUrl
	return nil
}

func (f *URLRepositoryTest) NextCount(ctx context.Context) (uint64, error) {
	return f.nextVal, nil
}

func (f *URLRepositoryTest) RedirectShortURL(ctx context.Context, shortURL string) (string, error) {
	return "", serviceErrors.ErrNotFound
}

func (f *URLRepositoryTest) UpdateCounter(ctx context.Context, shortURL string) error {
	return nil
}

func (f *URLRepositoryTest) GetCount(ctx context.Context, shortURL string) (int, error) {
	return 0, serviceErrors.ErrNotFound
}

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func() *URLRepositoryTest
		url     string
		wantURL string
		wantErr error
	}{
		{
			name: "новая ссылка",
			setup: func() *URLRepositoryTest {
				r := NewURLRepositoryTest()
				r.nextVal = 100000000001
				return r
			},
			url:     "https://example.com",
			wantURL: "http://localhost:8090/1l9Zo9p",
		},
		{
			name: "ссылка уже есть",
			setup: func() *URLRepositoryTest {
				r := NewURLRepositoryTest()
				r.shortByOriginal["https://example.com"] = "1l9Zo9p"
				return r
			},
			url:     "https://example.com",
			wantURL: "http://localhost:8090/1l9Zo9p",
		},
		{
			name: "невалидный URL",
			setup: func() *URLRepositoryTest {
				return NewURLRepositoryTest()
			},
			url:     "not-a-url",
			wantErr: serviceErrors.ErrBadRequest,
		},
		{
			name: "пустой URL",
			setup: func() *URLRepositoryTest {
				return NewURLRepositoryTest()
			},
			url:     "",
			wantErr: serviceErrors.ErrBadRequest,
		},
		{
			name: "GetShortURL ErrInternal",
			setup: func() *URLRepositoryTest {
				r := NewURLRepositoryTest()
				r.failGet = true
				return r
			},
			url:     "https://example.com",
			wantErr: serviceErrors.ErrInternal,
		},
		{
			name: "Insert ErrInternal",
			setup: func() *URLRepositoryTest {
				r := NewURLRepositoryTest()
				r.failInsert = true
				r.nextVal = 1
				return r
			},
			url:     "https://example.com",
			wantErr: serviceErrors.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.setup()
			svc := urlShort.NewReduceService(repo)

			got, err := svc.Create(context.Background(), tt.url)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("Create(%q) = %q, nil; want error %v", tt.url, got, tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Create(%q) err = %v, want %v", tt.url, err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Create(%q) неожиданная ошибка: %v", tt.url, err)
			}
			if got != tt.wantURL {
				t.Errorf("Create(%q) = %q, want %q", tt.url, got, tt.wantURL)
			}
		})
	}

}

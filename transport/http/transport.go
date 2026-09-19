package transport

import (
	"test/core/urlShort"
	"test/middleware"

	server "github.com/Elizabethppppp/tcp_server"
)

type Transport struct {
	tr *urlShort.ReduceService
}

func NewTransport(tr *urlShort.ReduceService) *Transport {
	return &Transport{
		tr: tr,
	}
}

func (t *Transport) Handler() *server.Mux {

	mux := server.NewMux()

	mux.Handle("post /short", middleware.LoggerMiddleware(t.CreateShortURL))
	mux.Handle("get /{short}", middleware.LoggerMiddleware(t.RedirectHandler))
	mux.Handle("GET /count/{short}", middleware.LoggerMiddleware(t.CountShortURL))

	return mux

}

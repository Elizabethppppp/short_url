package transport

import (
	"test/core/url_short"
	"test/middleware"

	server "github.com/Elizabethppppp/tcp_server"
)

type Transport struct {
	tr *url_short.ReduceService
}

func NewTransport(tr *url_short.ReduceService) *Transport {
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

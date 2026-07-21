package main

import (
	"context"
	"log"

	server "github.com/Elizabethppppp/tcp_server"
	"github.com/jackc/pgx/v5"
)

func main() {
	dsn := "postgres://postgres:hello1elephant@localhost:5432/URLstore?sslmode=disable"

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close(context.Background())

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("%v", err)
	}

	store := NewURLstore(db)

	mux := server.NewMux()
	mux.Handle("/short", store.CreateShortURL)
	mux.Handle("/{short}", store.RedirectHandler)
	//mux.Handle("/count/{short}", store.CountShortURL)

	if err := server.Listen(":8090", mux); err != nil {
		log.Fatal(err)
	}
}

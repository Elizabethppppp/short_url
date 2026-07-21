package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	server "github.com/Elizabethppppp/tcp_server"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	DB struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"db"`
}

func main() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	json.Unmarshal(file, &cfg)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)

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
	mux.Handle("/count/{short}", store.CountShortURL)

	if err := server.Listen(":8090", mux); err != nil {
		log.Fatal(err)
	}
}

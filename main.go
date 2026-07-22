package main

import (
	"encoding/json"
	"log"
	"os"
	db2 "test/db"

	server "github.com/Elizabethppppp/tcp_server"
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

	dbConn, err := db2.Connect(db2.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}

	defer dbConn.Close()

	store := NewURLstore(dbConn)

	mux := server.NewMux()
	mux.Handle("/short", store.CreateShortURL)
	mux.Handle("/{short}", store.RedirectHandler)
	mux.Handle("/count/{short}", store.CountShortURL)

	if err := server.Listen(":8090", mux); err != nil {
		log.Fatal(err)
	}
}

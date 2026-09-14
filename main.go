package main

import (
	"test/config"
	db2 "test/db"
	"test/logger"

	server "github.com/Elizabethppppp/tcp_server"
)

func main() {

	cfg, err := config.Load("config.yaml")
	if err != nil {
		panic(err)
	}

	if err := logger.Init(logger.Config{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
	}); err != nil {
		panic(err)
	}
	logger.Info("Config successfully parsed", "dbHost", cfg.DB.Host, "dbName", cfg.DB.DBName)

	logger.Debug("Database conection", "host", cfg.DB.Host, "port", cfg.DB.Port, "dbName", cfg.DB.DBName)

	dbConn, err := db2.Connect(db2.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
	})
	if err != nil {
		logger.Fatal("Fail connect database", err)
	}

	defer dbConn.Close()

	logger.Info("Connection successfully established", "host", cfg.DB.Host, "port", cfg.DB.Port)

	store := NewURLstore(dbConn)

	mux := server.NewMux()
	mux.Handle("post /short", LoggerMiddleware(store.CreateShortURL))
	mux.Handle("get /{short}", LoggerMiddleware(store.RedirectHandler))
	mux.Handle("GET /count/{short}", LoggerMiddleware(store.CountShortURL))

	logger.Info("Routes registered successfully")

	if err := server.Listen(":8090", mux); err != nil {
		logger.Fatal("Fail listen", err)
	}
}

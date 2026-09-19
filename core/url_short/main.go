package main

import (
	"test/config"
	db2 "test/db"
	"test/handler"
	"test/logger"
	"test/middleware"
	"test/pgService"
	"test/reduceService"

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
	logger.Info("Config successfully parsed", "dbHost", cfg.DB.Host, "dbName", cfg.DB.DBName, "dbSchema", cfg.DB.Schema)

	logger.Debug("Database conection", "host", cfg.DB.Host, "port", cfg.DB.Port, "dbName", cfg.DB.DBName)

	dbConn, err := db2.Connect(cfg.DB)
	if err != nil {
		logger.Fatal("Fail connect database", err)
	}

	defer dbConn.Close()

	logger.Info("Connection successfully established", "host", cfg.DB.Host, "port", cfg.DB.Port)

	store := handler.NewURLstore(reduceService.NewReduceService(pgService.NewPgService(dbConn)))

	mux := server.NewMux()
	mux.Handle("post /short", middleware.LoggerMiddleware(store.CreateShortURL))
	mux.Handle("get /{short}", middleware.LoggerMiddleware(store.RedirectHandler))
	mux.Handle("GET /count/{short}", middleware.LoggerMiddleware(store.CountShortURL))

	logger.Info("Routes registered successfully")

	if err := server.Listen(cfg.Server.Addr, mux); err != nil {
		logger.Fatal("Fail listen", err)
	}
}

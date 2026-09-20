package main

import (
	"log"
	"test/config"
	"test/core/urlShort"
	db2 "test/db"
	"test/logger"
	"test/pgService"
	transport "test/transport/http"

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

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	}()

	logger.Info("Connection successfully established", "host", cfg.DB.Host, "port", cfg.DB.Port)
	pg := pgService.NewPgService(dbConn)
	svc := urlShort.NewReduceService(pg)
	tp := transport.NewTransport(svc)

	handler := tp.Handler()
	logger.Info("Routes registered successfully")

	if err := server.Listen(cfg.Server.Addr, handler); err != nil {
		logger.Fatal("Fail listen", err)
	}
}

package main

import (
	"encoding/json"
	"os"
	db2 "test/db"
	"test/logger"

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
	Logger struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	} `json:"logger"`
}

func main() {

	if err := logger.Init(logger.Config{
		Level:  "info",
		Format: "json",
	}); err != nil {
		logger.Error("Fail init logger", err)
	}
	logger.Info("Starting")

	file, err := os.ReadFile("config.json")
	if err != nil {
		logger.Fatal("Fail read config.json", err)
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		logger.Fatal("Fail parse config.json", err)
	}

	logger.Info("Config successfully parsed", "dbHost", cfg.DB.Host, "dbName", cfg.DB.DBName)

	if cfg.Logger.Level != "" {
		err := logger.Init(logger.Config{
			Level:  cfg.Logger.Level,
			Format: cfg.Logger.Format,
		})
		if err != nil {
			logger.Error("Fail init logger", "error", err, "level", cfg.Logger.Level, "format", cfg.Logger.Format)
		} else {
			logger.Info("Logger successfully initialized", "level", cfg.Logger.Level, "format", cfg.Logger.Format)
		}
	}

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

package main

import (
	"flag"
	"log"

	"test/config"
	db2 "test/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configPath := flag.String("config", "config.yaml", "путь к config.yaml")
	migrationsDir := flag.String("dir", "db/migrations", "путь к папке с миграциями")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	dbConn, err := db2.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer dbConn.Close()

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{
		MigrationsTable: "schema_migrations",
		DatabaseName:    cfg.DB.DBName,
	})
	if err != nil {
		log.Fatalf("migrate driver: %v", err)
	}

	sourceURL := "file://" + *migrationsDir

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}

	log.Println("migrations applied successfully")
}

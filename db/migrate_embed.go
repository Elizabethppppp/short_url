package db

import "embed"

//go:embed migrations/*.sql
var MigrateFS embed.FS

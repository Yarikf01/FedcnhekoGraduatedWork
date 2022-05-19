package repo

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	log "github.com/Yarikf01/graduatedwork/api/utils"
)

func (db *DB) Migrate(ctx context.Context, path string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	logger := log.FromContext(ctx).With("mod", "DB migration")

	sourceURL := "file:///" + filepath.Join(dir, path, "migrations")
	logger.Infof("Looking for migration scripts in: %s\n", sourceURL)

	connStr := db.pool.Config().ConnString()

	m, err := migrate.New(sourceURL, connStr)
	if err != nil {
		return err
	}
	m.Log = &migrationLogger{logger: logger}

	if err = m.Up(); err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// helpers

type migrationLogger struct {
	logger *zap.SugaredLogger
}

func (l migrationLogger) Printf(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func (l migrationLogger) Verbose() bool {
	return true
}

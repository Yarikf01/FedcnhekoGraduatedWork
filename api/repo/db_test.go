package repo_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"

	"github.com/Yarikf01/graduatedwork/api/repo"
)

var (
	dbName    = "recon-test"
	pgxClient *pgx.Conn
	db        *repo.DB
)

// Setup & Cleanup resources for tests
func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	ctx := context.Background()

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("kartoza/postgis", "12.0", []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbName})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	connString := fmt.Sprintf("postgres://postgres:secret@%s/%s?sslmode=disable", resource.GetHostPort("5432/tcp"), dbName)

	db, err = repo.NewDB(ctx, connString, 2)
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if pgxClient, err = pgx.Connect(ctx, connString); err != nil {
		log.Fatalf("Could not connect with pgx: %s", err)
	}

	// create tables
	if err = db.Migrate(ctx, "../.."); err != nil {
		log.Fatalf("Migration failed: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestStorageConnection(t *testing.T) {
	err := pgxClient.Ping(context.Background())
	assert.NoError(t, err)

	err = db.Ping()
	assert.NoError(t, err)
}

func TestBulkInsertSQL(t *testing.T) {
	assert.Equal(t, "INSERT INTO table(id, action, description)"+
		" VALUES ($1, $2, $3),($4, $5, $6),($7, $8, $9),($10, $11, $12),($13, $14, $15)"+
		",($16, $17, $18),($19, $20, $21),($22, $23, $24),($25, $26, $27),($28, $29, $30)",
		repo.BulkInsertSQL("INSERT INTO table(id, action, description) VALUES %s", 3, 10))
}

func TestIsRowsNotFound(t *testing.T) {
	assert.True(t, repo.IsRowsNotFound(errors.New("no rows in result set")))
	assert.False(t, repo.IsRowsNotFound(errors.New("SQL error")))
}

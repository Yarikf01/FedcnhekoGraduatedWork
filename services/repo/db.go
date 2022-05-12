package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgxutil"

	"github.com/Yarikf01/graduatedwork/services/utils"
)

const (
	connectRetries = 10
	pingTimeout    = 3 * time.Second
)

var (
	ErrNotFound           = errors.New("record not found")
	NoRowsFound           = errors.New("no rows in result set")
	ErrPhoneAlreadyExists = errors.New("phone already exists")
)

const DefaultQueryLimit = 20

type TxRunner interface {
	RunWithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, connString string, maxConn int32) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url, %w", err)
	}

	poolConfig.ConnConfig.PreferSimpleProtocol = true
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
	}

	poolConfig.MaxConns = maxConn
	poolConfig.MinConns = 2

	poolConfig.ConnConfig.Logger = log.PgxLogAdapter{}
	poolConfig.ConnConfig.LogLevel = pgx.LogLevelWarn

	var pool *pgxpool.Pool
	operation := func() error {
		pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
		return err
	}

	err = backoff.Retry(operation, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), connectRetries))
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool, %w", err)
	}

	return &DB{pool: pool}, nil
}

func (db *DB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	_, err := db.pool.Exec(ctx, `select 1;`)
	return err
}

func (db *DB) Close() {
	db.pool.Close()
}

func BulkInsertSQL(insert string, argsPerRow, rows int) string {
	qMarks := make([]string, 0, argsPerRow)
	for i := 0; i < argsPerRow; i++ {
		qMarks = append(qMarks, "?")
	}
	rowValues := "(" + strings.Join(qMarks, ", ") + ")"
	// Combine the base SQL string and N value strings
	values := make([]string, 0, rows)
	for i := 0; i < rows; i++ {
		values = append(values, rowValues)
	}
	allValues := strings.Join(values, ",")
	insert = fmt.Sprintf(insert, allValues)
	// Convert all of the "?" to "$1", "$2", "$3", etc.
	numArgs := strings.Count(insert, "?")
	insert = strings.ReplaceAll(insert, "?", "$%v")
	numbers := make([]interface{}, 0, rows)
	for i := 1; i <= numArgs; i++ {
		numbers = append(numbers, strconv.Itoa(i))
	}
	return fmt.Sprintf(insert, numbers...)
}

func IsRowsNotFound(err error) bool {
	return err.Error() == NoRowsFound.Error()
}

// helpers

func (db *DB) selectInt(ctx context.Context, sql string) (int, error) {
	n, err := pgxutil.SelectInt64(ctx, db.doer(ctx), sql)
	return int(n), err
}

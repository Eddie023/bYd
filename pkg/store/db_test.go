package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/eddie023/byd/pkg/store"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/stretchr/testify/require"
)

type DBFixture struct {
	DB       *store.DB
	dbname   string
	hostname string
	conn     *sql.DB
}

func NewDBFixture(t require.TestingT) *DBFixture {
	dbname := "byd"

	hostname := os.Getenv("DB_HOSTNAME")
	if hostname == "" {
		hostname = "localhost"
	}

	conn, err := sql.Open("pgx", newPgx("postgres", hostname))
	require.NoError(t, err)
	defer conn.Close()
	require.NoError(t, conn.Ping())

	_, err = conn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", dbname))
	require.NoError(t, err)

	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
	require.NoError(t, err)

	conn, err = sql.Open("pgx", newPgx(dbname, hostname))
	require.NoError(t, err)

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance("file://../../migrations", "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)

	db, err := store.NewDB(context.Background(), newPgx(dbname, hostname))
	require.NoError(t, err)

	// read a seeds file
	seeds, err := os.ReadFile("../../migrations/seeds/insert_fakes.sql")
	require.NoError(t, err)

	_, err = conn.Exec(string(seeds))
	require.NoError(t, err)

	return &DBFixture{
		DB:       db,
		dbname:   dbname,
		hostname: hostname,
		conn:     conn,
	}
}

func newPgx(dbname, hostname string) string {
	host := "db"
	if hostname != "" {
		host = hostname
	}

	return fmt.Sprintf("postgres://root:postgres@%s:5432/%s?sslmode=disable", host, dbname)
}

func TestMigrationDown(t *testing.T) {
	f := NewDBFixture(t)

	conn, err := sql.Open("pgx", newPgx(f.dbname, f.hostname))
	require.NoError(t, err)

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance("file://../../migrations", "postgres", driver)
	require.NoError(t, err)

	err = m.Down()
	require.NoError(t, err)
}

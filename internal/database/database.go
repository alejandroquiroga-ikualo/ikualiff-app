package database

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ikualo.com/ikualiff/internal"
)

func CreateTable(sql string) (pgconn.CommandTag, error) {
	dbpool := getPool()
	defer dbpool.Close()

	return dbpool.Exec(context.Background(), sql)
}

func Exec(sql string, arguments ...any) error {
	dbpool := getPool()
	defer dbpool.Close()

	_, err := dbpool.Exec(context.Background(), sql, arguments...)
	return err
}

func QueryRow[T interface{}](sql string, args ...any) *T {
	dbpool := getPool()
	defer dbpool.Close()

	var result T
	err := dbpool.QueryRow(context.Background(), sql, args...).Scan(&result)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		log.Fatalf("Query row failed: %v", err)
	}

	return &result
}

func getPool() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), internal.GetEnv()["DATABASE_URL"])
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	return dbpool
}

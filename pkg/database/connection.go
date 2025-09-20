package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context,conn string) *sql.DB {
	db , err := sql.Open("postgres", conn)

	if err != nil {
		panic(err)
	}

	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	return db
}
package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*sqlx.DB, error) {
	connStr := "user=admin password=root dbname=go-chat sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Database) Close() {
	d.db.Close()
}

// func (d *Database) GetDB() *sql.DB {
// 	return d.db
// }

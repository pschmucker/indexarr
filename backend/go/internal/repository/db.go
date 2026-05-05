package repository

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dbPath string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Read and execute schema
	schema, err := os.ReadFile("internal/repository/schema.sql")
	if err != nil {
		return nil, err
	}

	if _, err := sqlDB.Exec(string(schema)); err != nil {
		return nil, err
	}

	db = sqlDB
	return sqlDB, nil
}

func GetDB() *sql.DB {
	return db
}

package repository

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrDBConnection = errors.New("DB connection failed")
	ErrDBMigration  = errors.New("DB migrations failed")
)

package database

import "database/sql"

type SQLConnector interface {
	Open(driverName, dsn string) (*sql.DB, error)
	Ping(db *sql.DB) error
}

type DefaultConnector struct{}

func (DefaultConnector) Open(driverName, dsn string) (*sql.DB, error) {
	return sql.Open(driverName, dsn)
}

func (DefaultConnector) Ping(db *sql.DB) error {
	return db.Ping()
}

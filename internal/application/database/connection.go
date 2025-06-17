package database

import (
	"abt-dashboard-api/pkg/database"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Username          string
	Password          string
	Database          string
	Host              string
	Port              string
	ReadTimeout       string
	WriteTimeout      string
	ConnectionTimeout string
}

func NewDbConfig(
	userName string,
	password string,
	database string,
	host string,
	port string,
	readTimeout string,
	writeTimeout string,
	connectionTimeout string,
) *DBConfig {
	return &DBConfig{
		Username:          userName,
		Password:          password,
		Database:          database,
		Host:              host,
		Port:              port,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ConnectionTimeout: connectionTimeout,
	}
}

func (d *DBConfig) NewDatabaseConnection(connector database.SQLConnector) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&readTimeout=%s&writeTimeout=%s&timeout=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.ReadTimeout,
		d.WriteTimeout,
		d.ConnectionTimeout,
	)

	db, err := connector.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(0)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)

	if err := connector.Ping(db); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

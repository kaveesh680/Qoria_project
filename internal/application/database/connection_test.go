package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConnector struct {
	openErr error
	pingErr error
	opened  bool
}

func (m *mockConnector) Open(driverName, dsn string) (*sql.DB, error) {
	m.opened = true
	if m.openErr != nil {
		return nil, m.openErr
	}
	return &sql.DB{}, nil
}

func (m *mockConnector) Ping(db *sql.DB) error {
	return m.pingErr
}

func TestNewDatabaseConnection_Success(t *testing.T) {
	cfg := NewDbConfig("user", "pass", "db", "localhost", "3306", "10s", "10s", "5s")
	mock := &mockConnector{}

	db, err := cfg.NewDatabaseConnection(mock)

	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.True(t, mock.opened)
}

func TestNewDatabaseConnection_OpenError(t *testing.T) {
	cfg := NewDbConfig("user", "pass", "db", "localhost", "3306", "10s", "10s", "5s")
	mock := &mockConnector{openErr: errors.New("open error")}

	db, err := cfg.NewDatabaseConnection(mock)

	assert.Nil(t, db)
	assert.EqualError(t, err, "failed to open database: open error")
}

func TestNewDatabaseConnection_PingError(t *testing.T) {
	cfg := NewDbConfig("user", "pass", "db", "localhost", "3306", "10s", "10s", "5s")
	mock := &mockConnector{pingErr: errors.New("ping error")}

	db, err := cfg.NewDatabaseConnection(mock)

	assert.Nil(t, db)
	assert.EqualError(t, err, "failed to ping database: ping error")
}

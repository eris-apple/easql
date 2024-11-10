package easql

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func setupMockService(client Client) (*Service, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(nil)

	config := &ConnectConfig{
		Client:   client,
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		User:     "user",
		Password: "password",
	}

	service := NewService(config, nil)
	service.SetConnection(mockDB)
	service.SetDatabase(gormDB)

	return service, mock
}

func TestService_Init_Postgres(t *testing.T) {
	service, mock := setupMockService(Postgres)

	mock.ExpectPing()

	err := service.Init()
	assert.NoError(t, err, "Service should initialize without error for Postgres")
	assert.NotNil(t, service.GetDatabase(), "Database should be initialized")
	assert.NotNil(t, service.GetConnect(), "Connection should be initialized")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestService_Init_MySQL(t *testing.T) {
	service, mock := setupMockService(MySQL)

	mock.ExpectPing()

	err := service.Init()
	assert.NoError(t, err, "Service should initialize without error for MySQL")
	assert.NotNil(t, service.GetDatabase(), "Database should be initialized")
	assert.NotNil(t, service.GetConnect(), "Connection should be initialized")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestService_Init_SQLite(t *testing.T) {
	service, mock := setupMockService(SQLite)

	mock.ExpectPing()

	err := service.Init()
	assert.NoError(t, err, "Service should initialize without error for SQLite")
	assert.NotNil(t, service.GetDatabase(), "Database should be initialized")
	assert.NotNil(t, service.GetConnect(), "Connection should be initialized")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestService_Init_UnsupportedClient(t *testing.T) {
	config := &ConnectConfig{
		Client:   "unsupported",
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		User:     "user",
		Password: "password",
	}

	service := NewService(config, nil)

	err := service.Init()
	assert.Error(t, err, "Service should return an error for unsupported client")
	assert.Contains(t, err.Error(), "sql client not support", "Error message should mention unsupported client")
}

func TestService_Disconnect_Success(t *testing.T) {
	service, mock := setupMockService(Postgres)

	mock.ExpectClose()

	err := service.Disconnect()
	assert.NoError(t, err, "Service should disconnect without error")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestService_Disconnect_Failure(t *testing.T) {
	service, mock := setupMockService(Postgres)

	mock.ExpectClose().WillReturnError(fmt.Errorf("close error"))

	err := service.Disconnect()
	assert.Error(t, err, "Service should return an error on failed disconnect")
	assert.Contains(t, err.Error(), "close error", "Error message should mention 'close error'")

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

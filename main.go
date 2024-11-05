package easql

import (
	"database/sql"
	"fmt"
	"github.com/eris-apple/ealogger"
	"github.com/eris-apple/eautils/url"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database = gorm.DB
type Connection = sql.DB

type Client = string

const (
	Postgres Client = "postgres"
	MySQL    Client = "mysql"
	SQLite   Client = "sqlite"
)

type ConnectConfig struct {
	Client   Client
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type Service struct {
	c *ConnectConfig
	l *ealogger.Logger

	conn     *Connection
	database *Database

	traceName string
}

func (s *Service) Init() error {
	s.l.DebugT(s.traceName, "Initializing sql service", s.c.Client)

	URL := url.NewURLConnectionString(s.c.Client, fmt.Sprintf("%s:%d", s.c.Host, s.c.Port), "", s.c.Database, s.c.User, s.c.Password)

	var dialect gorm.Dialector
	switch s.c.Client {
	case Postgres:
		dialect = postgres.Open(URL)
	case MySQL:
		dialect = mysql.Open(URL)
	case SQLite:
		dialect = sqlite.Open(URL)
	default:
		return fmt.Errorf("sql client not support: %s", s.c.Client)
	}

	db, err := gorm.Open(dialect, &gorm.Config{TranslateError: true})
	if err != nil {
		s.l.ErrorT(s.traceName, "Failed to connect to sql database", s.c.Client, err)
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		s.l.ErrorT(s.traceName, "Failed to connect to sql database", s.c.Client, err)
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		s.l.ErrorT(s.traceName, "Failed to ping sql database", s.c.Client, err)
		return err
	}

	s.database = db
	s.conn = sqlDB

	s.l.DebugT(s.traceName, "Postgres service initialized", s.c.Client)

	return nil
}

func (s *Service) Disconnect() error {
	if err := s.conn.Close(); err != nil {
		s.l.ErrorT(s.traceName, "Failed to close connection", s.c.Client, err)
	}

	return nil
}

func (s *Service) GetConnect() *Connection {
	return s.conn
}

func (s *Service) GetDatabase() *Database {
	return s.database
}

func NewService(c *ConnectConfig, l *ealogger.Logger, traceName string) *Service {
	return &Service{
		c: c,
		l: l,

		traceName: fmt.Sprintf("[%s_SQLService]", traceName),
	}
}

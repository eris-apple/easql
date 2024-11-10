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

const defaultTraceName = "[DEFAULT_SQLService]"

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

type ServiceConfig struct {
	IsLogging  bool
	Logger     *ealogger.Logger
	LoggerMode ealogger.Mode

	TraceName string
}

type Service struct {
	c *ConnectConfig
	l *ealogger.Logger

	Conn     *Connection
	Database *Database

	traceName string
}

func (s *Service) Init() error {
	s.l.DebugT(s.traceName, "Initializing sql service", s.c.Client)

	if s.Conn != nil || s.Database != nil {
		return nil
	}

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

	s.Database = db
	s.Conn = sqlDB

	s.l.DebugT(s.traceName, "Postgres service initialized", s.c.Client)

	return nil
}

func (s *Service) Disconnect() error {
	if s.Conn == nil {
		return fmt.Errorf("sql client not initialized")
	}

	if err := s.Conn.Close(); err != nil {
		s.l.ErrorT(s.traceName, "Failed to close connection", s.c.Client, err)
		return err
	}

	return nil
}

func (s *Service) GetConnect() *Connection {
	return s.Conn
}

func (s *Service) GetDatabase() *Database {
	return s.Database
}

func (s *Service) SetConnection(connection *Connection) {
	s.Conn = connection
}

func (s *Service) SetDatabase(db *Database) {
	s.Database = db
}

func (s *Service) SetTraceName(traceName string) {
	s.traceName = traceName
}

func (s *Service) SetLogger(logger *ealogger.Logger) {
	s.l = logger
}

func NewService(c *ConnectConfig, sc *ServiceConfig) *Service {
	if sc == nil {
		sc = &ServiceConfig{
			IsLogging: true,
			Logger:    ealogger.NewLoggerWithMode(ealogger.DevMode),
			TraceName: defaultTraceName,
		}
	} else {
		if sc.IsLogging && sc.Logger == nil {
			sc.Logger = ealogger.NewLoggerWithMode(ealogger.DevMode)
		}
	}

	return &Service{
		c: c,
		l: sc.Logger,

		traceName: sc.TraceName,
	}
}

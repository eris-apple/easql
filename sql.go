package easql

import (
	"database/sql"
	"fmt"
	"github.com/eris-apple/ealogger/ealogger"
	"github.com/eris-apple/eautils/url"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const defaultTraceName = "SQL_SERVICE"

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
	c         *ConnectConfig
	log       *ealogger.Entry
	IsLogging bool

	Conn     *Connection
	Database *Database

	traceName string
}

func (s *Service) Init() error {
	s.log.Debugf("Initializing sql service: %s", s.c.Client)

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
		s.log.Errorf("SQL client not support: %s", s.c.Client)
		return fmt.Errorf("sql client not support: %s", s.c.Client)
	}

	config := &gorm.Config{
		TranslateError: true,
	}

	if s.IsLogging {
		newLogger := logger.New(
			s.log,
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
				Colorful:                  false,
			},
		)
		config.Logger = newLogger
	}

	db, err := gorm.Open(dialect, config)
	if err != nil {
		s.log.Errorf("Failed to connect to sql database: %s, error: %s", s.c.Client, err)
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		s.log.Errorf("Failed to connect to sql database: %s, error: %s", s.c.Client, err)
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		s.log.Errorf("Failed to ping sql database: %s, error: %s", s.c.Client, err)
		return err
	}

	s.Database = db
	s.Conn = sqlDB

	s.log.Debugf("SQL Service initialized: %s", s.c.Client)

	return nil
}

func (s *Service) Disconnect() error {
	if s.Conn == nil {
		s.log.Errorf("SQL client not initialized: %s", s.c.Client)
		return fmt.Errorf("sql client not initialized")
	}

	if err := s.Conn.Close(); err != nil {
		s.log.Errorf("Failed to close connection: %s, error: %s", s.c.Client, err)
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

func (s *Service) SetLogger(logger *ealogger.Entry) {
	s.log = logger
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
		c:         c,
		log:       sc.Logger.WithName(sc.TraceName),
		IsLogging: sc.IsLogging,
	}
}

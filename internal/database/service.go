package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/aman/internal/logging"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service interface {
	GetDBConnection(ctx context.Context, dbName string) (*gorm.DB, error)
}

type service struct {
	logEnv           string
	maxDBConnections int
	maxConnLifetime  time.Duration
}

var (
	once     sync.Once
	instance *service
)

func New() Service {
	once.Do(func() {
		log := logging.Default()

		var cfg Config
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Errorw("error reading database config, using defaults", "err", err)
		}

		// ensure sensible default for duration if parsed as 0
		if cfg.MaxConnectionLifetimeMins == 0 {
			cfg.MaxConnectionLifetimeMins = 2 * time.Minute
		}

		instance = &service{
			logEnv:           cfg.LogEnv,
			maxDBConnections: cfg.DBMaxConn,
			maxConnLifetime:  cfg.MaxConnectionLifetimeMins,
		}
	})
	return instance
}


func (s *service) getTaskDBConfig() (*TaskDBConfig, error) {
	var cfg TaskDBConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read task db config: %w", err)
	}
	if cfg.DBURL == "" {
		return nil, errors.New("task DB URL is empty")
	}
	return &cfg, nil
}

func (s *service) GetDBConnection(ctx context.Context, dbName string) (*gorm.DB, error) {
	if dbName == "" {
		return nil, fmt.Errorf("dbName is empty")
	}

	// Read config (TASK_DB_URL / TASK_DB_READ_URL)
	cfg, err := s.getTaskDBConfig()
	if err != nil {
		return nil, fmt.Errorf("read task db config: %w", err)
	}

	// Compose write DSN for the specific database
	writeDSN, _, _, err := parseDSN(cfg.DBURL, dbName)
	if err != nil {
		return nil, fmt.Errorf("compose write dsn: %w", err)
	}

	// Open GORM DB
	gormDB, err := gorm.Open(postgres.Open(writeDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm open: %w", err)
	}

	// Configure underlying sql.DB
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("getting sql.DB from gorm: %w", err)
	}

	if s.maxDBConnections > 0 {
		sqlDB.SetMaxOpenConns(s.maxDBConnections)
		sqlDB.SetMaxIdleConns(s.maxDBConnections)
	}
	if s.maxConnLifetime > 0 {
		sqlDB.SetConnMaxLifetime(s.maxConnLifetime)
	}

	// Ping with context to ensure connectivity
	if err := sqlDB.PingContext(ctx); err != nil {
		// close the pool on error to avoid leaking resources
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return gormDB, nil
}

func parseDSN(rawURL, dbName string) (dsn, user, pass string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", "", err
	}
	if u.User == nil {
		return "", "", "", errors.New("DB URL missing user info")
	}
	user = u.User.Username()
	pass, _ = u.User.Password()
	if pass == "" {
		return "", "", "", errors.New("DB URL missing password")
	}
	if dbName != "" {
		u.Path = "/" + dbName
	}
	return u.String(), user, pass, nil
}

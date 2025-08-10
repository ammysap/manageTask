package database

import "time"

type Config struct {
	LogEnv                    string        `env:"LOG_ENV" env-default:"development"`
	DBMaxConn                 int           `env:"DB_MAX_CONN" env-default:"5"`
	MaxConnectionLifetimeMins time.Duration `env:"CONNECTION_LIFETIME" env-default:"2m"`
}

type TaskDBConfig struct {
	DBURL     string `env:"TASK_DB_URL" env-default:"postgres://user:password@localhost:5432/taskdb?sslmode=disable"`
	DBReadURL string `env:"TASK_DB_READ_URL"`
}
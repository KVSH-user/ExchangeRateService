// Package config is a package that provides a config.
package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config - ...
type Config struct {
	Env      string     `yaml:"env" env:"EXCHANGE_ENV"`
	Postgres PostgreSQL `yaml:"postgres" env:",inline"`
	REST     REST       `yaml:"rest" env:",inline"`
}

// PostgreSQL - ...
type PostgreSQL struct {
	Host              string        `yaml:"host" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_HOST"`
	MasterPort        int           `yaml:"master_port" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_MASTER_PORT"`
	DbName            string        `yaml:"db" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_DBNAME"`
	Login             string        `yaml:"login" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_USER"`
	LoginAdmin        string        `yaml:"login_admin" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_USER_ADMIN"`
	Password          string        `yaml:"password" env:"EXCHANGE_RATE_SERVICEG_POSTGRESQL_PASSWORD"`
	PasswordAdmin     string        `yaml:"password_admin" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_PASSWORD_ADMIN"`
	MaxConns          int32         `yaml:"max_conns" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_MAX_CONNS" env-default:"200"`
	MinConns          int32         `yaml:"min_conns" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_MIN_CONNS" env-default:"0"`
	MaxConnLifetime   time.Duration `yaml:"max_conn_lifetime" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_MAX_CONN_LIFETIME" env-default:"1h"`
	MaxConnIdleTime   time.Duration `yaml:"max_conn_idle_time" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_MAX_CONN_IDLE_TIME" env-default:"30m"`
	HealthcheckPeriod time.Duration `yaml:"healthcheck_period" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_HEALTHCHECK_PERIOD" env-default:"1m"`
	ConnectTimeout    time.Duration `yaml:"connect_timeout" env:"EXCHANGE_RATE_SERVICE_POSTGRESQL_CONNECT_TIMEOUT" env-default:"5s"`
}

// REST - ...
type REST struct {
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// Load - config load function
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config " + err.Error())
	}

	return &cfg
}

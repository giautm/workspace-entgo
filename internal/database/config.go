package database

import (
	"database/sql"
	"fmt"
)

// Config is the database configuration.
type Config struct {
	Name     string `env:"DB_NAME" json:",omitempty"`
	User     string `env:"DB_USER" json:",omitempty"`
	Host     string `env:"DB_HOST, default=localhost" json:",omitempty"`
	Port     string `env:"DB_PORT, default=3306" json:",omitempty"`
	Password string `env:"DB_PASSWORD" json:"-"` // ignored by zap's JSON formatter
}

// Database implements the Database interface.
func (c *Config) DatabaseConfig() *Config {
	return c
}

// ConnectionURL returns the connection URL.
func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

// OpenDB returns the database connection.
func (c *Config) OpenDB() *sql.DB {
	return sql.OpenDB(TraceConnector{
		DSN: c.ConnectionURL(),
	})
}

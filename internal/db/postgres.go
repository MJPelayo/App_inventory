// Package db handles PostgreSQL database connection and migrations
package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection pointer
var DB *sql.DB

// PostgresConfig holds all database connection parameters
type PostgresConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// GetPostgresConfig reads database configuration from environment variables
// Returns a config struct with all connection parameters
func GetPostgresConfig() *PostgresConfig {
    return &PostgresConfig{
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnv("DB_PORT", "5432"),
        User:     getEnv("DB_USER", "postgres"),
        Password: getEnv("DB_PASSWORD", ""),
        DBName:   getEnv("DB_NAME", "inventory_db"),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
    }
}

// ConnectionString builds the PostgreSQL connection string from config
// Format: postgres://user:password@host:port/dbname?sslmode=disable
func (c *PostgresConfig) ConnectionString() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// InitDB establishes connection to PostgreSQL and verifies it works
func InitDB() error {
    config := GetPostgresConfig()
    
    var err error
    DB, err = sql.Open("postgres", config.ConnectionString())
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    // Configure connection pool
    DB.SetMaxOpenConns(25)
    DB.SetMaxIdleConns(10)
    DB.SetConnMaxLifetime(5 * time.Minute)

    // Test the connection
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    log.Println("✅ PostgreSQL database connected successfully")
    return nil
}

// CloseDB closes the database connection (called on server shutdown)
func CloseDB() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}

// getEnv reads environment variable with fallback default value
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
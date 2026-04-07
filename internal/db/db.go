package db
// Package db handles database connection, migrations, and queries
package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3" // SQLite driver - no other dependencies needed
)

// DB is the global database connection pointer used by all handlers
var DB *sql.DB

// InitDB opens the database connection and runs all migration files
// Called once when the server starts up
func InitDB(dataSourceName string) error {
    var err error
    
    // Open SQLite database file (creates file if it doesn't exist)
    DB, err = sql.Open("sqlite3", dataSourceName)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    // Verify the database connection is working
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    log.Println("✅ Database connected successfully")

    // Run all SQL migration files to create tables and seed data
    if err = runMigrations(); err != nil {
        return fmt.Errorf("failed to run migrations: %w", err)
    }

    return nil
}

// runMigrations reads every .sql file in the migrations folder and executes it
// Files run in alphabetical order (001_, 002_, etc.) to ensure proper table creation order
func runMigrations() error {
    // Path to the migrations folder relative to project root
    migrationsDir := "internal/db/migrations"
    
    // Read all file names in the migrations directory
    files, err := os.ReadDir(migrationsDir)
    if err != nil {
        return fmt.Errorf("failed to read migrations directory: %w", err)
    }

    // Execute each SQL file in alphabetical order
    for _, file := range files {
        // Only process files with .sql extension
        if filepath.Ext(file.Name()) == ".sql" {
            // Read the entire SQL file content as a string
            path := filepath.Join(migrationsDir, file.Name())
            content, err := os.ReadFile(path)
            if err != nil {
                return fmt.Errorf("failed to read migration %s: %w", file.Name(), err)
            }

            // Execute the SQL commands against the database
            _, err = DB.Exec(string(content))
            if err != nil {
                return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
            }
            log.Printf("✅ Migration %s executed", file.Name())
        }
    }

    return nil
}

// CloseDB closes the database connection (called when server shuts down)
func CloseDB() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}
package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// RunMigrations executes all SQL migration files in order
// Files must be named: 001_name.up.sql, 002_name.up.sql, etc.
func RunMigrations() error {
	// Path to migrations folder
	migrationsDir := "internal/db/migrations"

	// Read all files in migrations directory
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Collect only .up.sql files
	var migrationFiles []string
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() && filepath.Ext(name) == ".sql" {
			migrationFiles = append(migrationFiles, name)
		}
	}

	// Sort files alphabetically (ensures 001, 002, 003 order)
	sort.Strings(migrationFiles)

	// Execute each migration file
	for _, fileName := range migrationFiles {
		path := filepath.Join(migrationsDir, fileName)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", fileName, err)
		}

		// Execute the SQL commands
		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}
		log.Printf("✅ Migration %s executed", fileName)
	}

	return nil
}

package migrate

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migrator handles database migrations
type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	// Get the directory where the executable is located
	exePath, err := os.Executable()
	if err != nil {
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	
	// Try to find migrations directory relative to executable
	migrationsDir := filepath.Join(exeDir, "db", "migrations")
	
	// If not found, try relative to current working directory
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		migrationsDir = filepath.Join(cwd, "db", "migrations")
	}
	
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Migrate runs all up migrations in order
func (m *Migrator) Migrate() error {
	// Create migrations table if not exists
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY NOT NULL,
			applied_at TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of applied migrations
	appliedVersions, err := m.getAppliedVersions()
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// Get all migration files
	migrationFiles, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Sort migrations by version
	sort.Strings(migrationFiles)

	// Apply pending migrations
	for _, file := range migrationFiles {
		version := extractVersion(file)
		if version == "" {
			continue
		}

		// Skip if already applied
		if appliedVersions[version] {
			continue
		}

		// Read and execute migration
		filePath := filepath.Join(m.migrationsDir, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		// Execute migration
		_, err = m.db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		// Record migration as applied
		_, err = m.db.Exec(`
			INSERT INTO schema_migrations (version, applied_at) VALUES (?, datetime('now'))
		`, version)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}
	}

	return nil
}

// getAppliedVersions returns a map of applied migration versions
func (m *Migrator) getAppliedVersions() (map[string]bool, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}

	return versions, rows.Err()
}

// getMigrationFiles returns a list of .up.sql migration files
func (m *Migrator) getMigrationFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(m.migrationsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".up.sql") {
			files = append(files, info.Name())
		}
		return nil
	})

	return files, err
}

// extractVersion extracts the version number from a migration filename
func extractVersion(filename string) string {
	if !strings.HasSuffix(filename, ".up.sql") {
		return ""
	}
	parts := strings.Split(filename, "_")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

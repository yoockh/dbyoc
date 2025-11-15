package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a single database migration.
type Migration struct {
	Version     int
	Description string
	Script      string
}

// MigrationManager manages the migrations.
type MigrationManager struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new MigrationManager.
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{
		db: db,
	}
}

// LoadMigrations loads migration scripts from the specified directory.
func (m *MigrationManager) LoadMigrations(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			version, description, err := parseMigrationFile(file.Name())
			if err != nil {
				return err
			}

			script, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			m.migrations = append(m.migrations, Migration{
				Version:     version,
				Description: description,
				Script:      string(script),
			})
		}
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	return nil
}

// ApplyMigration applies a migration to the database.
func (m *MigrationManager) ApplyMigration(migration Migration) error {
	_, err := m.db.Exec(migration.Script)
	if err != nil {
		return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
	}
	return nil
}

// parseMigrationFile extracts version and description from the migration file name.
func parseMigrationFile(filename string) (int, string, error) {
	parts := strings.SplitN(filename, "__", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid migration file name format: %s", filename)
	}

	version := parts[0]
	description := strings.TrimSuffix(parts[1], filepath.Ext(parts[1]))

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		return 0, "", fmt.Errorf("invalid migration version: %s", version)
	}

	return versionInt, description, nil
}

// GetMigrations returns the list of loaded migrations.
func (m *MigrationManager) GetMigrations() []Migration {
	return m.migrations
}
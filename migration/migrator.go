package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate(migrationsDir string) error {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			err := m.applyMigration(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", file.Name(), err)
			}
		}
	}
	return nil
}

func (m *Migrator) applyMigration(filePath string) error {
	migrationSQL, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", filePath, err)
	}

	_, err = m.db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	fmt.Printf("Migration %s applied successfully\n", filepath.Base(filePath))
	return nil
}

func (m *Migrator) Rollback(migrationsDir string) error {
	// Implement rollback logic if needed
	return nil
}

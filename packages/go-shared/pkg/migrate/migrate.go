package migrate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Version     string
	Description string
	SQL         string
}

func Run(db *sqlx.DB, migrationsDir string, dbName string) error {
	ensureMigrationsTable(db)

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var migrations []Migration
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		parts := strings.SplitN(f.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}
		content, err := os.ReadFile(filepath.Join(migrationsDir, f.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f.Name(), err)
		}
		migrations = append(migrations, Migration{
			Version:     parts[0],
			Description: strings.TrimSuffix(parts[1], ".sql"),
			SQL:         string(content),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, m := range migrations {
		if applied[m.Version] {
			continue
		}
		log.Printf("[migrate] applying %s_%s (%s)", m.Version, m.Description, dbName)
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for %s: %w", m.Version, err)
		}
		if _, err := tx.Exec(m.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", m.Version, err)
		}
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version, description) VALUES (?, ?)",
			m.Version, m.Description,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", m.Version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", m.Version, err)
		}
		log.Printf("[migrate] applied %s_%s OK", m.Version, m.Description)
	}

	return nil
}

func ensureMigrationsTable(db *sqlx.DB) {
	db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(10) PRIMARY KEY,
		description VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
}

func getAppliedMigrations(db *sqlx.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return make(map[string]bool), nil
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, nil
}

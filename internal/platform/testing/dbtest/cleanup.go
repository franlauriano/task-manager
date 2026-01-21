package dbtest

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
)

var cleanupMutex sync.Mutex

// CleanDatabase truncates all tables in the database using TRUNCATE CASCADE.
// It dynamically queries all tables in the public schema via information_schema.
// Useful for resetting database state between tests.
// Thread-safe: uses mutex to prevent concurrent cleanup operations.
func CleanDatabase(db *gorm.DB) error {
	cleanupMutex.Lock()
	defer cleanupMutex.Unlock()

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var tableNames []string
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := tx.Raw(query).Rows()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to iterate tables: %w", err)
	}

	if len(tableNames) == 0 {
		return tx.Commit().Error
	}

	truncateSQL := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE",
		quoteIdentifiers(tableNames))

	if err := tx.Exec(truncateSQL).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit cleanup transaction: %w", err)
	}

	return nil
}

// quoteIdentifiers formats a list of table names for use in SQL.
// Ensures names with special characters are properly quoted (PostgreSQL double-quote style).
func quoteIdentifiers(names []string) string {
	quoted := make([]string, len(names))
	for i, name := range names {
		quoted[i] = fmt.Sprintf(`"%s"`, name)
	}
	return strings.Join(quoted, ", ")
}

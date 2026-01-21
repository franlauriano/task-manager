package dbtest

import (
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// LoadFixtures loads one or more fixture SQL files from a directory and commits the transaction.
// If names are provided, only those files are loaded; otherwise all .sql files in the directory are loaded.
// Used for loading test fixtures (e.g. db/fixtures).
func LoadFixtures(db *gorm.DB, fixtureDir string, names ...string) error {
	absDir, err := filepath.Abs(fixtureDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute fixture path: %w", err)
	}

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

	if err := loadFixtures(tx, absDir, names...); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit fixtures transaction: %w", err)
	}

	return nil
}

// ResetWithFixtures clears all tables and loads the given fixture files.
func ResetWithFixtures(db *gorm.DB, fixtureDir string, names ...string) error {
	if err := CleanDatabase(db); err != nil {
		return err
	}
	return LoadFixtures(db, fixtureDir, names...)
}

// loadFixtures loads and executes SQL files from a directory.
// If names are provided, only loads those specific files. Otherwise, loads all .sql files.
func loadFixtures(tx *gorm.DB, directory string, names ...string) error {
	// If specific names provided, load only those
	if len(names) > 0 {
		for _, name := range names {
			filePath := filepath.Join(directory, name)
			if filepath.Ext(name) != ".sql" {
				filePath += ".sql"
			}
			if err := executeSQLFile(tx, filePath); err != nil {
				return err
			}
		}
		return nil
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		if err := executeSQLFile(tx, filePath); err != nil {
			return err
		}
	}

	return nil
}

// executeSQLFile reads and executes SQL from a file in the given transaction.
func executeSQLFile(tx *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filepath.Base(filePath), err)
	}

	if err := tx.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("failed to execute SQL from %s: %w", filepath.Base(filePath), err)
	}

	return nil
}

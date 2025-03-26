package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	drivePostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

type DBConnector interface {
	GetDB() *gorm.DB
	Close() error
}

func New(connection string) (*Database, error) {
	db, err := gorm.Open(drivePostgres.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("❌ can't connect to database: %w", err)
	}

	database := &Database{DB: db}

	if err := database.RunMigrations(); err != nil {
		return nil, fmt.Errorf("❌ migration error: %w", err)
	}

	return &Database{DB: db}, nil
}

func (p *Database) GetDB() *gorm.DB {
	return p.DB
}

func (p *Database) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func getMigrationPath() string {
	baseDir, _ := os.Getwd() // Lấy thư mục làm việc hiện tại (cmd/identity)
	migrationsPath := filepath.Join(baseDir, "..", "..", "pkg", "identity", "db", "migrations")

	return fmt.Sprintf("file://%s", migrationsPath)
}

func (p *Database) RunMigrations() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		getMigrationPath(),
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	version, dirty, _ := m.Version()
	if dirty {
		fmt.Println("⚠️ Previous migration encountered an error, please check!")
		return fmt.Errorf("migration is in an error state (dirty)")
	}

	fmt.Printf("📌 Current migration version: %d\n", version)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

package postgres

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
	Ping() error
}

func New(connection, serviceName string) (*Database, error) {
	db, err := gorm.Open(drivePostgres.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("❌ can't connect to database: %w", err)
	}

	database := &Database{DB: db}

	if err := database.RunMigrations(serviceName); err != nil {
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

func (p *Database) Ping() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func getMigrationPath(serviceName string) string {
	baseDir, _ := os.Getwd()
	migrationsPath := filepath.Join(baseDir, "..", "..", "pkg", serviceName, "db", "migrations")

	return fmt.Sprintf("file://%s", migrationsPath)
}

func (p *Database) RunMigrations(serviceName string) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		getMigrationPath(serviceName),
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if dirty {
		forceTo := int(version) - 1
		fmt.Printf("⚠️ Dirty version detected (%d). Forcing back to version %d\n", version, forceTo)

		if err := m.Force(forceTo); err != nil {
			return fmt.Errorf("❌ failed to force clean migration version: %w", err)
		}
	}

	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("❌ failed to get migration version: %w", err)
	}

	fmt.Printf("📌 Current migration version: %d\n", version)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

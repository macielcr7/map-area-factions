package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/macielcr7/map-area-factions/backend/internal/config"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DatabaseService struct {
	DB *gorm.DB
}

func NewDatabaseService(cfg *config.DatabaseConfig) (*DatabaseService, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	fmt.Printf("Database config: host=%s, port=%s, user=%s, dbname=%s\n",
		cfg.Host, cfg.Port, cfg.User, cfg.Name)

	// Configure GORM logger
	var gormLogger logger.Interface
	gormLogger = logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "public.",
			SingularTable: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &DatabaseService{DB: db}, nil
}

func (ds *DatabaseService) AutoMigrate() error {
	if err := ds.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("failed to ensure uuid-ossp extension: %w", err)
	}

	if err := ds.ensureCoreSchema(); err != nil {
		return err
	}

	const createSystemSettingsTable = `
	CREATE TABLE IF NOT EXISTS system_settings (
	    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	    data JSONB NOT NULL,
	    created_at TIMESTAMPTZ DEFAULT NOW(),
	    updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	`

	if err := ds.DB.Exec(createSystemSettingsTable).Error; err != nil {
		return fmt.Errorf("failed to ensure system_settings table: %w", err)
	}

	return nil
}

func (ds *DatabaseService) ensureCoreSchema() error {
	var hasUsersTable bool
	checkSQL := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema = CURRENT_SCHEMA()
			  AND table_name = ?
		);
	`

	if err := ds.DB.Raw(checkSQL, "users").Scan(&hasUsersTable).Error; err != nil {
		return fmt.Errorf("failed to verify users table: %w", err)
	}

	if !hasUsersTable {
		migrationPath := filepath.Join("migrations", "001_initial_schema.up.sql")
		sqlBytes, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationPath, err)
		}

		script := strings.TrimSpace(string(sqlBytes))
		if script == "" {
			return fmt.Errorf("migration file %s is empty", migrationPath)
		}

		if err := ds.DB.Exec(script).Error; err != nil {
			return fmt.Errorf("failed to run initial schema migration: %w", err)
		}
	}

	return ds.ensureSoftDeleteColumns()
}

func (ds *DatabaseService) ensureSoftDeleteColumns() error {
	tables := []struct {
		model  interface{}
		column string
	}{
		{&models.User{}, "deleted_at"},
		{&models.Faction{}, "deleted_at"},
		{&models.Region{}, "deleted_at"},
		{&models.Geometry{}, "deleted_at"},
		{&models.Report{}, "deleted_at"},
	}

	for _, table := range tables {
		if ds.DB.Migrator().HasColumn(table.model, table.column) {
			continue
		}
		if err := ds.DB.Migrator().AddColumn(table.model, table.column); err != nil {
			return fmt.Errorf("failed to add column %s: %w", table.column, err)
		}
	}

	return nil
}

func (ds *DatabaseService) Close() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (ds *DatabaseService) HealthCheck() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

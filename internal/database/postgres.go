package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	URL string
}

func NewPostgresConnection(cfg Config) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return gormDB, nil
}

package database

import (
	"fmt"
	"log"

	"github.com/zenkriztao/ayo-football-backend/internal/config"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase creates a new database connection based on configuration
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "postgres":
		var dsn string
		if cfg.Database.Password == "" {
			dsn = fmt.Sprintf(
				"postgres://%s@%s:%s/%s?sslmode=%s",
				cfg.Database.User,
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)
		} else {
			dsn = fmt.Sprintf(
				"postgres://%s:%s@%s:%s/%s?sslmode=%s",
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)
		}
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	logLevel := logger.Silent
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate entities
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// autoMigrate runs auto migration for all entities
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Team{},
		&entity.Player{},
		&entity.Match{},
		&entity.Goal{},
	)
}

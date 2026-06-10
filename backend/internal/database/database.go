package database

import (
	"fmt"
	"time"

	"github.com/example/ems/internal/config"
	authmodel "github.com/example/ems/internal/auth/model"
	empmodel "github.com/example/ems/internal/employee/model"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New opens a GORM connection to PostgreSQL and configures the connection pool.
func New(cfg config.DatabaseConfig, log *zap.Logger, production bool) (*gorm.DB, error) {
	gormLogLevel := logger.Info
	if production {
		gormLogLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(gormLogLevel),
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		return nil, fmt.Errorf("database: open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	log.Info("database connected",
		zap.String("host", cfg.Host),
		zap.String("db", cfg.Name),
	)
	return db, nil
}

// Migrate auto-migrates the schema. The uuid-ossp extension is enabled so that
// gen_random_uuid()/uuid_generate_v4() defaults work. For production prefer the
// versioned SQL migrations in /migrations over AutoMigrate.
func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("database: enable uuid-ossp: %w", err)
	}
	if err := db.AutoMigrate(&authmodel.User{}, &empmodel.Employee{}); err != nil {
		return fmt.Errorf("database: automigrate: %w", err)
	}
	return nil
}

// Close releases the underlying connection pool.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// Give in-flight queries a brief window before closing.
	time.Sleep(50 * time.Millisecond)
	return sqlDB.Close()
}

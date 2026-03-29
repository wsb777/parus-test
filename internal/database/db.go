package database

import (
	"fmt"
	"parus-test/internal/config"
	"parus-test/internal/database/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DatabasePG struct {
	DBUsername string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string

	logger *zap.Logger
}

func NewDatabasePG(cfg *config.Config, logger *zap.Logger) *DatabasePG {
	return &DatabasePG{
		DBUsername: cfg.DBUser,
		DBPassword: cfg.DBPassword,
		DBName:     cfg.DBName,
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,

		logger: logger,
	}
}

func (d *DatabasePG) Connect() (*gorm.DB, error) {
	d.logger.Info("connecting to database", zap.String("host", d.DBHost), zap.String("db", d.DBName))

	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC"
	dsn = fmt.Sprintf(dsn, d.DBHost, d.DBUsername, d.DBPassword, d.DBName, d.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		d.logger.Fatal("failed to connect to database", zap.Error(err))
		return nil, err
	}

	d.logger.Info("connected to database")
	return db, nil
}

func (d *DatabasePG) Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Group{}, &models.FileEntry{}, &models.FileVersion{}, &models.Token{})

	if err != nil {
		d.logger.Fatal("failed to run migrations", zap.Error(err))
		return err
	}

	d.logger.Info("migrations completed")
	return nil
}

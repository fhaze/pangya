package database

import (
	"fmt"
	"os"
	"pangya/src/internal/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TIME_ZONE"),
		os.Getenv("DB_SSL_MODE"),
	)
	logger.Log.Info("connecting to database", zap.String("dsn", dsn))
	gormDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	logger.Log.Info("successfully conected to database")
	db = gormDb
	return nil
}

func Get() *gorm.DB {
	return db
}

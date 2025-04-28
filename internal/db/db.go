package db

import (
	"bank-app-backend/internal/entities"
	"bank-app-backend/internal/lib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(path string) (*gorm.DB, error) {
	var err error

	db, err = gorm.Open(postgres.Open(path), &gorm.Config{})
	if err != nil {
		logger.Log.Fatal("Could not connect to database", zap.Error(err))
	}

	if err := db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.Account{}); err != nil {
		logger.Log.Fatal("Could not migrate database", zap.Error(err))
	}

	return db, nil
}

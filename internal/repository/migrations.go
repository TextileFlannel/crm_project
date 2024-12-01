package repository

import (
	"http-server/danilkovalev/internal/models"

	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
		&models.Integration{},
		&models.Contact{},
	)
}

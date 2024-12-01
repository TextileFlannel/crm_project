package repository

import (
	"fmt"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

func NewMySQLDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	fmt.Println("Connected to the database successfully")
	return db, nil
}

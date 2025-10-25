package utils

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// Hardcoded database credentials
	DB_USER := "root"
	DB_PASS := "golfantara"
	DB_HOST := "localhost"
	DB_NAME := "case_study_2"
	DB_PORT := "3306"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	return db
}

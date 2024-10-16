package config

import (
	"codebase-api/internal/domain"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// Optional: migrasikan schema jika perlu
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&domain.User{})
	return db
}

func GetDBDSN() string {
	return os.Getenv("DB_DSN")
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

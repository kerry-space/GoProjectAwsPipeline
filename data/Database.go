package data

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func openMySql(server, database, username, password string, port int) *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, server, port, database)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}

func Init(file, server, database, username, password string, port int) {
	if len(file) == 0 {
		DB = openMySql(server, database, username, password, port)
	} else {
		DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
	}
	// AutoMigrate for both Car and User models
	if err := DB.AutoMigrate(&Car{}, &User{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	seedCars()
	seedUsers()
}

func seedCars() {
	var count int64
	DB.Model(&Car{}).Count(&count)
	if count == 0 {
		DB.Create(&Car{1, "GMC Yukon denali", "2024", "black"})
		DB.Create(&Car{2, "Tesla X", "2025", "blackOut"})
		DB.Create(&Car{3, "Chevrolet Cameron", "2025", "blackOut"})
	}
}

func seedUsers() {
	var count int64
	DB.Model(&User{}).Count(&count)
	if count == 0 {
		// Seed with a default user or setup an admin user etc.
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("yourSecurePassword"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("failed to hash password: %v", err)
		}
		DB.Create(&User{Username: "admin", Password: string(hashedPassword)})
	}
}

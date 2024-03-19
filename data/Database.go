package data

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func openMySql(server, database, username, password string, port int) *gorm.DB {
	var url string
	url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, server, port, database)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func Init(file, server, database, username, password string, port int) {
	if len(file) == 0 {
		DB = openMySql(server, database, username, password, port)
	} else {
		DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
	}
	DB.AutoMigrate(&Car{})
	var antal int64
	DB.Model(&Car{}).Count(&antal) // Seed
	if antal == 0 {
		DB.Create(&Car{1, "GMC Yukon denali", "2024", "black"})
		DB.Create(&Car{2, "Tesla X", "2025", "blackOut"})
		DB.Create(&Car{3, "Chevrolet Cameron", "2025", "blackOut"})
	}

}

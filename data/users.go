package data

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:255;uniqueIndex"` // Explicitly define size and unique index
	Password string `gorm:"size:255"`             // Define size if necessary
}

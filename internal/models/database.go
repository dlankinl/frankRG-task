package models

import "github.com/jinzhu/gorm"

var DB *gorm.DB

func SetDatabase(db *gorm.DB) {
	DB = db
}

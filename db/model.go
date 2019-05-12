package db

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func Conn() *gorm.DB {
	return db
}

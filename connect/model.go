package connect

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

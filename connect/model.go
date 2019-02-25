package model

import (
	"github.com/jinzhu/gorm"
)

// главная конфигурация
type configTyp struct {
	Mysql      mysql
	Postgresql postgresql
}

// конфигурация поджгружаемая из файла конфигурации
type mysql struct {
	Host     string // протокол, хост и порт подключения
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type postgresql struct {
	Host     string // Хост базы данных (localhost - по умолчанию)
	Port     int64  // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

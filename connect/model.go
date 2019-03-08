package connect

import (
	"github.com/jinzhu/gorm"
)

// главная конфигурация
type Config struct {
	Mysql      mysql      `yaml:"Mysql"`
	Postgresql postgresql `yaml:"Postgresql"`
}

// конфигурация поджгружаемая из файла конфигурации
type mysql struct {
	Host     string `yaml:"Host"`     // протокол, хост и порт подключения
	Name     string `yaml:"Name"`     // Имя базы данных
	Login    string `yaml:"Login"`    // Логин к базе данных
	Password string `yaml:"Password"` // Пароль к базе данных
	Charset  string `yaml:"Charset"`  // Кодировка данных (utf-8 - по умолчанию)
}

type postgresql struct {
	Host     string `yaml:"Host"`     // Хост базы данных (localhost - по умолчанию)
	Port     int64  `yaml:"Port"`     // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string `yaml:"Name"`     // Имя базы данных
	Login    string `yaml:"Login"`    // Логин к базе данных
	Password string `yaml:"Password"` // Пароль к базе данных
	Charset  string `yaml:"Charset"`  // Кодировка данных (utf-8 - по умолчанию)
	Ssl      string `yaml:"Ssl"`      // Ssl
}

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

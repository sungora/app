package conf

import (
	"os"
	"time"
)

func init() {
	dir, _ := os.Getwd()
	sep := string(os.PathSeparator)
	ConfigDir = dir + sep + "config" + sep
}

type ConfigMain struct {
	Name     string
	TimeZone string
	DriverDB string // Драйвер DB
	Mode     string // Режим работы приложения
	Host     string
	Port     int
}

type ConfigMysql struct {
	Host     string // протокол, хост и порт подключения
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type ConfigPostgresql struct {
	Host     string // Хост базы данных (localhost - по умолчанию)
	Port     int64  // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

var ConfigDir string
var TimeLocation *time.Location

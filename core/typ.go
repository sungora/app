package core

import (
	"time"

	"github.com/sungora/app/v2/lg"
	"github.com/sungora/app/v2/workflow"
)

type config struct {
	Main       configMain
	Mysql      configMysql
	Postgresql configPostgresql
	Log        lg.Config
	Workflow   workflow.Config
}

type configMain struct {
	TimeZone       string        // Временная зона
	DriverDB       string        // Драйвер DB
	SessionTimeout time.Duration // Время жизни сессии в секундах
	Host           string        // Хост веб сервера
	Port           int           // Порт веб сервера
	Mode           string        // Режим работы приложения
	Static         string        // Папка для статических данными (css, js, img, etc...)
	Template       string        // Папка для шаблонов
}

type configMysql struct {
	Host     string // протокол, хост и порт подключения
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type configPostgresql struct {
	Host     string // Хост базы данных (localhost - по умолчанию)
	Port     int64  // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

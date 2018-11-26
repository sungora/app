package core

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ConfigTyp struct {
	Name           string        // Техническое название приложения
	TimeZone       string        // Временная зона
	DriverDB       string        // Драйвер DB
	SessionTimeout time.Duration // Время жизни сессии в секундах
	Host           string        // Хост веб сервера
	Port           int           // Порт веб сервера
	Mode           string        // Режим работы приложения
	Static         string        // Папка для статических данными (css, js, img, etc...)
	Template       string        // Папка для шаблонов
}

var Config *ConfigTyp
var DB *gorm.DB

package servhttp

import (
	"net/http"
	"time"
)

// компонент
type Component struct {
	Server    *http.Server  // сервер HTTP
	chControl chan struct{} // управление ожиданием завершения работы сервера
}

// конфигурация
type Config struct {
	Proto          string        `yaml:"Proto"`          // Server Proto
	Host           string        `yaml:"Host"`           // Server Host
	Port           int           `yaml:"Port"`           // Server Port
	ReadTimeout    time.Duration `yaml:"ReadTimeout"`    // Время ожидания web запроса в секундах, по истечении которого соединение сбрасывается
	WriteTimeout   time.Duration `yaml:"WriteTimeout"`   // Время ожидания окончания передачи ответа в секундах, по истечении которого соединение сбрасывается
	IdleTimeout    time.Duration `yaml:"IdleTimeout"`    // Время ожидания следующего запроса
	MaxHeaderBytes int           `yaml:"MaxHeaderBytes"` // Максимальный размер заголовка получаемого от браузера клиента в байтах
	Cors           Cors          `yaml:"Cors"`           // cors headers
	Proxy          string        `yaml:"Proxy"`          // format: "http://login:pass@bproxy.msk.mts.ru:3131"
}

type Cors struct {
	Use              bool     `yaml:"Use"`
	AllowedOrigins   []string `yaml:"AllowedOrigins"`
	AllowedMethods   []string `yaml:"AllowedMethods"`
	AllowedHeaders   []string `yaml:"AllowedHeaders"`
	ExposedHeaders   []string `yaml:"ExposedHeaders"`
	AllowCredentials bool     `yaml:"AllowCredentials"`
	MaxAge           int      `yaml:"MaxAge"`
}

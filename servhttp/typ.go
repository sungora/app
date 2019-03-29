package servhttp

// конфигурация
type Config struct {
	Proto          string `yaml:"Proto"`          // Server Proto
	Host           string `yaml:"Host"`           // Server Host
	Port           int    `yaml:"Port"`           // Server Port
	ReadTimeout    int    `yaml:"ReadTimeout"`    // Время ожидания web запроса в секундах, по истечении которого соединение сбрасывается
	WriteTimeout   int    `yaml:"WriteTimeout"`   // Время ожидания окончания передачи ответа в секундах, по истечении которого соединение сбрасывается
	IdleTimeout    int    `yaml:"IdleTimeout"`    // Время ожидания следующего запроса
	MaxHeaderBytes int    `yaml:"MaxHeaderBytes"` // Максимальный размер заголовка получаемого от браузера клиента в байтах
	Cors           Cors   `yaml:"Cors"`           // cors headers
}

type Cors struct {
	IsUse            bool     `yaml:"IsUse"`
	AllowedOrigins   []string `yaml:"AllowedOrigins"`
	AllowedMethods   []string `yaml:"AllowedMethods"`
	AllowedHeaders   []string `yaml:"AllowedHeaders"`
	ExposedHeaders   []string `yaml:"ExposedHeaders"`
	AllowCredentials bool     `yaml:"AllowCredentials"`
	MaxAge           int      `yaml:"MaxAge"`
}

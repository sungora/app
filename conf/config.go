package conf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	NameApp     string
	DisplayName string
	Description string
	TimeZone    string
	DriverDB    string // Драйвер DB
	Mode        string // Режим работы приложения
	AutoRestart bool
	Server      Server
	Mysql       Mysql
	Postgresql  Postgresql
	Log         Log
	Workflow    Workflow
}

type Server struct {
	Host string
	Port int
}

type Mysql struct {
	Host     string // протокол, хост и порт подключения
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type Postgresql struct {
	Host     string // Хост базы данных (localhost - по умолчанию)
	Port     int64  // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
}

type Log struct {
	ServiceName string
	Info        bool
	Notice      bool
	Warning     bool
	Error       bool
	Critical    bool
	Fatal       bool
	Debug       bool
	Traces      bool
	OutStd      bool
	OutFile     bool
	OutFilePath string
	OutHttp     bool
	OutHttpUrl  string // url куда отправляются логи
}

type Workflow struct {
	Isworkflow bool
	LimitCh    int // Лимит канала задач
	LimitPool  int // Лимит пула (количество воркеров)
}

var conf *Config

// ReadToml Функция читает конфигурационный файл в формате toml. Отдельный конфиг не связанный с beego.
func GetConfig() (*Config, error) {
	if conf != nil {
		return conf, nil
	}
	var (
		configFile string
		nameApp    string
		dir        string
		err        error
	)
	sep := string(os.PathSeparator)
	if dir, err = os.Getwd(); err == nil {
		if ext := filepath.Ext(os.Args[0]); ext != "" {
			sl := strings.Split(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
			nameApp = sl[0]
		} else {
			nameApp = filepath.Base(os.Args[0])
		}
	}
	configFile = dir + sep + "config" + sep + "main.toml"
	if _, err := toml.DecodeFile(configFile, &conf); err != nil {
		return nil, err
	}
	if conf.TimeZone == "" {
		conf.TimeZone = "Europe/Moscow"
	}
	conf.NameApp = nameApp
	return conf, nil
}

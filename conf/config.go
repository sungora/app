package conf

import (
	"os"
)

func init() {
	dir, _ := os.Getwd()
	sep := string(os.PathSeparator)
	ConfigDir = dir + sep + "config" + sep
}

var ConfigDir string

var Main *Config

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



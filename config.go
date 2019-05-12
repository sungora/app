package app

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

const (
	extToml = "toml"
	extYaml = "yaml"
)

// конфигурация
type Config struct {
	SessionTimeout time.Duration  `yaml:"SessionTimeout"` //
	TimeZone       string         `yaml:"TimeZone"`       //
	Domain         string         `yaml:"Domain"`         //
	Mode           string         `yaml:"Mode"`           //
	DirWork        string         `yaml:"DirWork"`        //
	ServiceName    string         `yaml:"ServiceName"`    // Техническое название приложения
	TimeLocation   *time.Location ``                      // Временная зона
	Version        string         `yaml:"Version"`        // Версия приложения
}

// конфигурация
var config *Config

// SConfigearchPath поиск конфигурации
func ConfigSearchPath(serviceName string) (path, ext string) {
	if serviceName == "" {
		serviceName = filepath.Base(os.Args[0])
		serviceName = strings.Split(serviceName, filepath.Ext(serviceName))[0]
	}
	//
	sep := string(os.PathSeparator)
	path = filepath.Dir(filepath.Dir(os.Args[0]))
	if path == "." {
		path, _ = os.Getwd()
		path = filepath.Dir(path)
	}
	path += sep + "config" + sep + serviceName
	//
	if inf, err := os.Stat(path + extYaml); err == nil {
		if inf.Mode().IsRegular() == true {
			return path, extYaml
		}
	}
	if inf, err := os.Stat(path + extToml); err == nil {
		if inf.Mode().IsRegular() == true {
			return path, extToml
		}
	}
	return
}

// ConfigLoad загрузка конфигурации
func ConfigLoad(path string, cfg interface{}) (err error) {
	var data []byte
	l := strings.SplitAfter(path, ".")
	ext := l[len(l)-1]
	switch ext {
	case extToml:
		_, err = toml.DecodeFile(path, cfg);
	case extYaml:
		if data, err = ioutil.ReadFile(path); err != nil {
			return
		}
		if err = yaml.Unmarshal(data, cfg); err != nil {
			return
		}
	default:
		return errors.New("undefined config: " + path)
	}
	return
}

// ConfigSetDefault инициализация дефолтовыми значениями
func ConfigSetDefault(cfg *Config) {
	if cfg == nil {
		cfg = &Config{}
	}
	// временная зона
	if cfg.TimeZone != "" {
		cfg.TimeZone = "Europe/Moscow"
	}
	if loc, err := time.LoadLocation(cfg.TimeZone); err == nil {
		cfg.TimeLocation = loc
	} else {
		cfg.TimeLocation = time.UTC
	}
	// режим работы приложения
	if cfg.Mode == "" {
		cfg.Mode = "dev"
	}
	// техническое имя приложения
	if cfg.ServiceName == "" {
		if ext := filepath.Ext(os.Args[0]); ext != "" {
			sl := strings.Split(filepath.Base(os.Args[0]), ext)
			cfg.ServiceName = sl[0]
		} else {
			cfg.ServiceName = filepath.Base(os.Args[0])
		}
	}
	// пути
	sep := string(os.PathSeparator)
	if cfg.DirWork == "" {
		cfg.DirWork, _ = os.Getwd()
		sl := strings.Split(cfg.DirWork, sep)
		if sl[len(sl)-1] == "bin" {
			sl = sl[:len(sl)-1]
		}
		cfg.DirWork = strings.Join(sl, sep)
	}
	// сессия
	if cfg.SessionTimeout == 0 {
		cfg.SessionTimeout = 86400
	}
	//
	if cfg.Domain == "" {
		cfg.Domain = "localhost"
	}
	if cfg.Version == "" {
		cfg.Version = "1.0.0"
	}
	//
	config = cfg
}

func GetConfig() *Config {
	if config == nil {
		ConfigSetDefault(nil)
	}
	return config
}

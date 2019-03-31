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

	"github.com/sungora/app/session"
)

const (
	extToml = "toml"
	extYaml = "yaml"
)

// конфигурация
type Config struct {
	SessionTimeout time.Duration  `yaml:"SessionTimeout"` //
	TimeZone       string         `yaml:"TimeZone"`       //
	Mode           string         `yaml:"Mode"`           //
	DirWork        string         `yaml:"DirWork"`        //
	DirConfig      string         `yaml:"DirConfig"`      //
	DirLog         string         `yaml:"DirLog"`         //
	DirWww         string         `yaml:"DirWww"`         //
	ServiceName    string         `yaml:"ServiceName"`    // Техническое название приложения
	TimeLocation   *time.Location ``                      // Временная зона
	Version        string         `yaml:"Version"`        // Версия приложения
}

// конфигурация
var Cfg *Config

func configuration(cfg *Config) {
	Cfg = cfg
	// временная зона
	if Cfg.TimeZone != "" {
		Cfg.TimeZone = "Europe/Moscow"
	}
	if loc, err := time.LoadLocation(Cfg.TimeZone); err == nil {
		Cfg.TimeLocation = loc
	} else {
		Cfg.TimeLocation = time.UTC
	}
	// режим работы приложения
	if Cfg.Mode == "" {
		Cfg.Mode = "dev"
	}
	// техническое имя приложения
	if Cfg.ServiceName == "" {
		if ext := filepath.Ext(os.Args[0]); ext != "" {
			sl := strings.Split(filepath.Base(os.Args[0]), ext)
			Cfg.ServiceName = sl[0]
		} else {
			Cfg.ServiceName = filepath.Base(os.Args[0])
		}
	}
	// пути
	sep := string(os.PathSeparator)
	if Cfg.DirWork == "" {
		Cfg.DirWork, _ = os.Getwd()
		// Cfg.DirWork = filepath.Dir(filepath.Dir(os.Args[0]))
		// if Cfg.DirWork == "." {
		// 	Cfg.DirWork, _ = os.Getwd()
		// 	Cfg.DirWork = filepath.Dir(Cfg.DirWork)
		// }
	}
	if Cfg.DirConfig == "" {
		Cfg.DirConfig = Cfg.DirWork + sep + "config"
	}
	if Cfg.DirLog == "" {
		Cfg.DirLog = Cfg.DirWork + sep + "log"
	}
	if Cfg.DirWww == "" {
		Cfg.DirWww = Cfg.DirWork + sep + "www"
	}
	// сессия
	Cfg.SessionTimeout *= time.Second
	if 0 < Cfg.SessionTimeout {
		session.SessionGC(Cfg.SessionTimeout)
	}
}

// SearchConfigPath поиск конфигурации
func SearchConfigPath(serviceName string) (path, ext string) {
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

// LoadConfig загрузка конфигурации
func LoadConfig(path string, cfg interface{}) (err error) {
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

package core

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

func GetConfig() *Config {
	return config
}

// LoadConfigYaml загрузка конфигурации в формате yaml
func LoadConfigYaml(path string, cfg interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return
	}
	return
}

// LoadConfigToml загрузка конфигурации в формате toml
func LoadConfigToml(path string, cfg interface{}) (err error) {
	// читаем конфигурацию
	_, err = toml.DecodeFile(path, cfg);
	return
}

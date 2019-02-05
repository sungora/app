package config

import (
	"os"

	"github.com/BurntSushi/toml"

	"github.com/sungora/app/core"
)

// init регистрация компонента в приложении
func init() {
	component = new(componentTyp)
	core.ComponentReg(component)
}

// компонент
type componentTyp struct {
}

var (
	config    *configMain   // конфигурация
	component *componentTyp // компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *core.ConfigRoot) (err error) {
	sep := string(os.PathSeparator)
	config = new(configMain)

	// читаем конфигурацию
	path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, config); err != nil {
		return
	}

	// Назанчение конфигураций
	ServiceAccess = config.SampleConfig

	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {
	return
}

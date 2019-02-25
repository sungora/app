package model

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"

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
	config    *configTyp    // конфигурация
	component *componentTyp // компонент
)

// Init инициализация компонента в приложении
func (comp *componentTyp) Init(cfg *core.ConfigRoot) (err error) {

	sep := string(os.PathSeparator)
	config = new(configTyp)

	// читаем конфигурацию
	path := cfg.DirConfig + sep + cfg.ServiceName + ".toml"
	if _, err = toml.DecodeFile(path, config); err != nil {
		return
	}

	if config.Mysql.Host != "" {
		if db, err = gorm.Open("mysql", fmt.Sprintf(
			"%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
			config.Mysql.Login,
			config.Mysql.Password,
			config.Mysql.Host,
			config.Mysql.Name,
			config.Mysql.Charset,
		)); err != nil {
			return
		}
	} else if config.Postgresql.Host != "" {
		if db, err = gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s",
			config.Postgresql.Host,
			config.Postgresql.Port,
			config.Postgresql.Login,
			config.Postgresql.Name,
			config.Postgresql.Password,
		)); err != nil {
			return
		}
	}
	return
}

// Start запуск компонента в работу
func (comp *componentTyp) Start() (err error) {
	return
}

// Stop завершение работы компонента
func (comp *componentTyp) Stop() (err error) {

	if db != nil {
		err = db.Close()
	}

	return
}

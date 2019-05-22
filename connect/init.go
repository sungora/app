// Deprecated
// Use db
package connect

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	config    *Config    // конфигурация
	component *Component // компонент
)

// Init инициализация компонента в приложении
func Init(cfg *Config, IsLog bool) (com *Component, err error) {
	config = cfg
	component = new(Component)
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
		if IsLog == true {
			db = db.Debug()
		}
	} else if config.Postgresql.Host != "" {
		if db, err = gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			config.Postgresql.Host,
			config.Postgresql.Port,
			config.Postgresql.Login,
			config.Postgresql.Name,
			config.Postgresql.Password,
			config.Postgresql.Ssl,
		)); err != nil {
			return
		}
		if IsLog == true {
			db = db.Debug()
		}
	}
	return component, nil
}

// Start запуск компонента в работу
func (comp *Component) Start() (err error) {
	return
}

// Stop завершение работы компонента
func (comp *Component) Stop() (err error) {

	if db != nil {
		err = db.Close()
	}

	return
}

func GetConfig() *Config {
	return config
}

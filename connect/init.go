package connect

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// компонент
type Component struct {
}

var (
	config    *Config    // конфигурация
	component *Component // компонент
)

// Init инициализация компонента в приложении
func Init(cfg *Config) (com *Component, err error) {
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
	} else if config.Postgresql.Host != "" {
		if db, err = gorm.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s charset=%s",
			config.Postgresql.Host,
			config.Postgresql.Port,
			config.Postgresql.Login,
			config.Postgresql.Name,
			config.Postgresql.Password,
			config.Postgresql.Ssl,
			config.Postgresql.Charset,
		)); err != nil {
			return
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

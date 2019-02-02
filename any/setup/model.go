package setup

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"accounter/model"
)

func initModel(config *Config) (err error) {
	switch config.App.DriverDB {
	case "mysql":
		if model.DB, err = gorm.Open("mysql", fmt.Sprintf(
			"%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
			config.Mysql.Login,
			config.Mysql.Password,
			config.Mysql.Host,
			config.Mysql.Name,
			config.Mysql.Charset,
		)); err != nil {
			return
		}
	case "postgresql":
		if model.DB, err = gorm.Open("postgres", fmt.Sprintf(
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

func waitModel() {
	if model.DB != nil {
		model.DB.Close()
	}
}

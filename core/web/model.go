package web

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/sungora/app.v1/conf"
)

var TimeLocation *time.Location
var DB *gorm.DB

func Start(conf *conf.Config) (err error) {
	if loc, err := time.LoadLocation(conf.TimeZone); err == nil {
		TimeLocation = loc
	} else {
		TimeLocation = time.UTC
	}
	switch conf.DriverDB {
	case "mysql":
		if err = StartMysql(conf.Mysql, conf.TimeZone); err != nil {
			return
		}
	case "postgresql":
		if err = StartPostgresql(conf.Postgresql, conf.TimeZone); err != nil {
			return
		}
	}
	return
}

func StartMysql(conf conf.Mysql, timeZone string) (err error) {
	if loc, err := time.LoadLocation(timeZone); err == nil {
		TimeLocation = loc
	} else {
		TimeLocation = time.UTC
	}
	DB, err = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@%s/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
		conf.Login, conf.Password, conf.Host, conf.Name, conf.Charset,
	))
	if err != nil {
		return err
	}
	return nil
}

func StartPostgresql(conf conf.Postgresql, timeZone string) (err error) {
	if loc, err := time.LoadLocation(timeZone); err == nil {
		TimeLocation = loc
	} else {
		TimeLocation = time.UTC
	}
	DB, err = gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		conf.Host, conf.Port, conf.Login, conf.Name, conf.Password,
	))
	if err != nil {
		return err
	}
	return nil
}

func Wait() {
	if DB != nil {
		DB.Close()
	}
}

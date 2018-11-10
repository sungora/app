// TODO реализовать именованые конфиги
package database

import (
	"gopkg.in/kshamiev/sungora.v1/database/face"
	"gopkg.in/kshamiev/sungora.v1/database/mysql"
	"gopkg.in/kshamiev/sungora.v1/lg"
)

type DbFace face.DbFace
type ArFace face.ArFace

var driver string

// Доступные БД для использования

func NewDb() DbFace {
	switch driver {
	case "mysql":
		return mysql.NewDb()
	}
	lg.Fatal(168, driver)
	return nil
}

func NewAr() ArFace {
	switch driver {
	case "mysql":
		return mysql.NewAr()
	}
	return nil
}

// CheckConnect Проверка настроек, конфигарций и соединений с БД
func CheckConnect(dr string) (err error) {
	driver = dr
	switch driver {
	case "mysql":
		return mysql.CheckConnect()
	}
	return
}

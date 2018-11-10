package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"gopkg.in/kshamiev/sungora.v1/database/mysql/confmysql"
)

// Стек конфигураций БД
var cfg confmysql.Mysql
var timeLocation *time.Location

func Start(c confmysql.Mysql, timeZone string) {
	if c.Type == `` {
		c.Type = `tcp`
	}
	if c.Host == `` {
		c.Host = `localhost`
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	if c.Charset == `` {
		c.Charset = `UTF-8`
	}
	if c.TimeOut == 0 {
		c.TimeOut = 5
	}
	if c.CntConn == 0 {
		c.CntConn = 50
	}
	// Инициализация временной зоны
	if loc, err := time.LoadLocation(timeZone); err == nil {
		timeLocation = loc
	} else {
		timeLocation = time.UTC
	}
	cfg = c

	// контроль коннектов
	go func() {
		for {
			for key := range conn {
				t := conn[key].time.Add(time.Second * time.Duration(cfg.TimeOut))
				if conn[key].free == true && 0 < time.Now().In(timeLocation).Sub(t) {
					conn[key].Connect.Close()
					delete(conn, key)
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()
}

// CheckConnect Проверка настроек, конфигарций и соединений с БД
func CheckConnect() (err error) {
	var db mysql.Conn
	d := cfg
	if d.Type == `tcp` {
		if "" == d.Host || 0 == d.Port || "" == d.Login || "" == d.Name {
			return errors.New("Настройки Mysql указаны не полностью для tcp")
		}
		db = mysql.New("tcp", "", fmt.Sprintf("%s:%d", d.Host, d.Port), d.Login, d.Password)
	} else {
		if "" == d.Socket || "" == d.Host || 0 == d.Port || "" == d.Login || "" == d.Name {
			return errors.New("Настройки Mysql указаны не полностью для socket")
		}
		db = mysql.New("unix", "", d.Socket, d.Login, d.Password)
	}
	err = db.Connect()
	if err != nil {
		return
	}
	err = db.Use(d.Name)
	if err != nil {
		return
	}
	db.Close()
	return
}

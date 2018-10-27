package core

import (
	"gopkg.in/sungora/app.v1/conf"
	"time"
)

func init() {
	// Запуск чистки старых сессий по таймауту
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			for i, s := range session {
				if 100 < time.Now().Sub(s.t) {
					delete(session, i)
				}
			}
		}
	}()
}

var session = make(map[string]*Session)

type Session struct {
	t    time.Time
	data map[string]interface{}
}

func GetSession(token string) *Session {
	if elm, ok := session[token]; ok {
		return elm
	}
	session[token] = new(Session)
	session[token].t = time.Now().In(conf.TimeLocation)
	session[token].data = make(map[string]interface{})
	return session[token]
}

func (self *Session) Get(index string) interface{} {
	if elm, ok := self.data[index]; ok {
		return elm
	}
	return nil
}

func (self *Session) Set(index string, value interface{}) {
	self.data[index] = value
}
func (self *Session) Del(index string) {
	delete(self.data, index)
}

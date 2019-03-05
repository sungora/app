package core

import (
	"time"
)

// sessionGC Запуск чистки старых сессий по таймауту
func sessionGC() {
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			for i, s := range session {
				if Config.SessionTimeout < time.Now().In(Config.TimeLocation).Sub(s.t) {
					delete(session, i)
				}
			}
		}
	}()
}

var session = make(map[string]*sessionTyp)

type sessionTyp struct {
	t    time.Time
	data map[string]interface{}
}

// GetSession Получение сессии
func GetSession(token string) *sessionTyp {
	if elm, ok := session[token]; ok {
		elm.t = time.Now().In(Config.TimeLocation)
		return elm
	}
	session[token] = new(sessionTyp)
	session[token].t = time.Now().In(Config.TimeLocation)
	session[token].data = make(map[string]interface{})
	return session[token]
}

// Get получение данных сессии
func (s *sessionTyp) Get(index string) interface{} {
	if _, ok := s.data[index]; ok {
		return s.data[index]
	}
	return nil
}

// Set сохранение данных в сессии
func (s *sessionTyp) Set(index string, value interface{}) {
	s.data[index] = value
}

// Del удаление данных из сессии
func (s *sessionTyp) Del(index string) {
	delete(s.data, index)
}

package core

import (
	"time"

	"gopkg.in/sungora/app.v1/tool"
)

// SessionGC Запуск чистки старых сессий по таймауту
func SessionGC() {
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			for i, s := range session {
				if Config.Main.SessionTimeout < time.Now().In(tool.TimeLocation).Sub(s.t) {
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
		elm.t = time.Now().In(tool.TimeLocation)
		return elm
	}
	session[token] = new(Session)
	session[token].t = time.Now().In(tool.TimeLocation)
	session[token].data = make(map[string]interface{})
	return session[token]
}

func (s *Session) Get(index string) interface{} {
	if elm, ok := s.data[index]; ok {
		return elm
	}
	return nil
}

func (s *Session) Set(index string, value interface{}) {
	s.data[index] = value
}

func (s *Session) Del(index string) {
	delete(s.data, index)
}

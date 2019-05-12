package session

import (
	"time"
)

// SessionGC Запуск чистки старых сессий по таймауту
func SessionGC(td time.Duration) {
	go func() {
		for {
			time.Sleep(time.Minute)
			for i, s := range sessions {
				if td < time.Now().Sub(s.t) {
					delete(sessions, i)
				}
			}
		}
	}()
}

var sessions = make(map[string]*Session)

type Session struct {
	t    time.Time
	data map[string]interface{}
}

// GetSession Получение сессии
func GetSession(token string) *Session {
	if elm, ok := sessions[token]; ok {
		elm.t = time.Now()
		return elm
	}
	sessions[token] = new(Session)
	sessions[token].t = time.Now()
	sessions[token].data = make(map[string]interface{})
	return sessions[token]
}

// Get получение данных сессии
func (s *Session) Get(index string) interface{} {
	if _, ok := s.data[index]; ok {
		return s.data[index]
	}
	return nil
}

// Set сохранение данных в сессии
func (s *Session) Set(index string, value interface{}) {
	s.data[index] = value
}

// Del удаление данных из сессии
func (s *Session) Del(index string) {
	delete(s.data, index)
}

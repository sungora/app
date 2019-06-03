package servws

import (
	"fmt"
	"time"
)

// интерфейс сообщения и его обработчиков
type Message interface {
	HookStartClient()
	HookSendMessage()
	String() string
}

// пример сообщения чата и его обработчики
type MessageSample struct {
	ClientID   string    `json:"client_id"`
	Author     string    `json:"author"`
	Body       string    `json:"body"`
	BodyBinary string    `json:"body_binary"`
	CreatedAt  time.Time `json:"created_at"`
}

// HookStartClient обработка входящего сообщения
func (m *MessageSample) HookStartClient() () {
	fmt.Println("WS hook start client")
}

// HookSend обработка входящего сообщения
func (m *MessageSample) HookSendMessage() {
	fmt.Println("WS hook send message: ", m)
}

func (m *MessageSample) String() string {
	return m.ClientID + " - " + m.Author + " - " + m.Body
}

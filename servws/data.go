package servws

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// интерфейс сообщения и его обработчиков
type Message interface {
	HookStartClient(ws *websocket.Conn, cnt int)
	HookGetMessage(ws *websocket.Conn, cnt int)
	HookSendMessage(ws *websocket.Conn, cnt int)
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
func (m *MessageSample) HookStartClient(ws *websocket.Conn, cnt int) () {
	fmt.Println("WS hook start client")
}

// HookSend обработка отправляемого сообщения на сервер
// (для других пользователей)
func (m *MessageSample) HookSendMessage(ws *websocket.Conn, cnt int) {
	fmt.Println("WS hook send message: ", m)
}

// HookSend обработка входящего сообщения от сервера
// (от других пользователей)
func (m *MessageSample) HookGetMessage(ws *websocket.Conn, cnt int) {
	fmt.Println("WS hook get message: ", m)
}

// String удобоваримый вывод значения структуры
func (m *MessageSample) String() string {
	return m.ClientID + " - " + m.Author + " - " + m.Body
}

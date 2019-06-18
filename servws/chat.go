package servws

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// управление удалением пустых (покинутых) чатов
func ControlBusChat() {
	for {
		select {
		case <-time.After(time.Second * 10):
			for i := range busChat {
				if len(busChat[i].clients) == 0 {
					delete(busChat, i)
				}
			}
		}
	}
}

// шина чатов
var busChat = map[string]*BusChat{}

type registerClientChat struct {
	ws      *websocket.Conn
	handler Message
}

// чат
type BusChat struct {
	register  chan *registerClientChat    // канал регистрации нового клиента
	broadcast chan Message                // канал рассылки сообщений клиентам
	clients   map[*websocket.Conn]Message // массив всех клиентов чата
}

// InitChat инициализация чата по условному идентификатору
func InitChat(chatID string) *BusChat {
	if _, ok := busChat[chatID]; ok == false {
		b := &BusChat{
			register:  make(chan *registerClientChat),
			broadcast: make(chan Message),
			clients:   make(map[*websocket.Conn]Message),
		}
		go b.control()
		busChat[chatID] = b
		fmt.Println("WS new chat: " + chatID)
	} else {
		fmt.Println("WS get chat: " + chatID)
	}
	return busChat[chatID]
}

// start управление чатом
func (b *BusChat) control() {
	for {
		select {
		// каждому зарегистрированному клиенту шлем сообщение
		case message := <-b.broadcast:
			for client, handler := range b.clients {
				// если достучаться до клиента не удалось, то удаляем его
				if _, err := client.NextWriter(websocket.PingMessage); err != nil {
					delete(b.clients, client)
					continue
				}
				// hook handler get other client message
				handler.HookGetMessage(client)
				if err := client.WriteJSON(message); err != nil {
					fmt.Println("WS error send message")
				}
			}
		// регистрируем новго клиента
		case client := <-b.register:
			b.clients[client.ws] = client.handler
			// hook handler new client
			client.handler.HookStartClient(client.ws)
		}
	}
}

// StartClient регистрация и старт работы клиента
func (b *BusChat) StartClient(ws *websocket.Conn, msg Message) {
	if msg == nil {
		msg = &MessageSample{}
	}
	b.register <- &registerClientChat{ws, msg}
	for {
		// var msg Message
		if err := ws.ReadJSON(msg); err != nil {
			delete(b.clients, ws)
			return
		} else {
			// hook handler send owher client message
			msg.HookSendMessage(ws)
			// посылаем всем подключенным пользователям
			b.broadcast <- msg
		}
	}
}

package servws

import (
	"fmt"
	"io"
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

type registerClientChat struct {
	ws      *websocket.Conn
	handler Message
}

// чат
type BusChat struct {
	register  chan registerClientChat     // канал регистрации нового клиента
	broadcast chan Message                // канал рассылки сообщений клиентам
	clients   map[*websocket.Conn]Message // массив всех клиентов чата
}

// шина чатов
var busChat = map[string]*BusChat{}

// InitChat инициализация чата по условному идентификатору
func InitChat(chatID string) *BusChat {
	if _, ok := busChat[chatID]; ok == false {
		b := &BusChat{
			register:  make(chan registerClientChat),
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
		// проверка соединений с клиентами
		case <-time.After(time.Second * 50):
			for client := range b.clients {
				// если достучаться до клиента не удалось, то удаляем его
				if _, err := client.NextWriter(websocket.PingMessage); err != nil {
					delete(b.clients, client)
					continue
				}
			}
		// каждому зарегистрированному клиенту шлем сообщение
		case message := <-b.broadcast:
			for client, handler := range b.clients {
				// hook handler get other client message
				handler.HookGetMessage(client, len(b.clients))
				if err := client.WriteJSON(message); err != nil {
					fmt.Println("WS error send message")
				}
			}
		// регистрируем новго клиента
		case client := <-b.register:
			b.clients[client.ws] = client.handler
			// hook handler new client
			client.handler.HookStartClient(client.ws, len(b.clients))
		}
	}
}

// StartClient регистрация и старт работы клиента
func (b *BusChat) StartClient(ws *websocket.Conn, msg Message) {
	if msg == nil {
		msg = &MessageSample{}
	}
	b.register <- registerClientChat{ws, msg}
	for {
		// var msg Message
		err := ws.ReadJSON(msg)
		if err == io.ErrUnexpectedEOF {
			delete(b.clients, ws)
			return
		} else if err != nil {
			fmt.Println("WS error parsing message")
		} else {
			// hook handler send owher client message
			msg.HookSendMessage(ws, len(b.clients))
			// посылаем всем подключенным пользователям
			b.broadcast <- msg
		}
	}
}

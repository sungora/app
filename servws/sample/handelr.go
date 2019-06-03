package sample

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"

	"github.com/sungora/app/request"
	"github.com/sungora/app/servws"
)

func init() {
	go servws.ControlBusChat() // контроль за чат комнатами (удаление пустых - покинутых)
}

// ChatStart
// http://localhost:8080/chat/index.html
func ChatStart(w http.ResponseWriter, r *http.Request) {
	//
	// Здесь можно расположить проверку можно ли запустить чат или нет
	//
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		request.NewIn(w, r).JsonError(-1, err.Error())
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			fmt.Println("WS close connect error:" + err.Error())
		} else {
			fmt.Println("WS close connect ok")
		}
	}()
	bus := servws.InitChat(chi.URLParam(r, "id"))
	bus.StartClient(ws, &servws.MessageSample{})
}

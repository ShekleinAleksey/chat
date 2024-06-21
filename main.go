package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ConnectUser struct {
	websocket *websocket.Conn
	ClientIP  string
}

func newConnectUser(ws *websocket.Conn, ClientIP string) *ConnectUser {
	return &ConnectUser{
		websocket: ws,
		ClientIP:  ClientIP,
	}
}

func main() {
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func init() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/ws", WebSocketHandler)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("template/index.html")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var users = make(map[ConnectUser]int)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, _ := upgrader.Upgrade(w, r, nil)

	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	log.Println("Client connected:", ws.RemoteAddr().String())

	var socketClient *ConnectUser = newConnectUser(ws, ws.RemoteAddr().String())
	users[*socketClient] = 0
	log.Println("Number client connected ...", len(users))

	messageType, message, err := ws.ReadMessage()
	if err != nil {
		log.Println("Ws disconnect waiting", err.Error())
		delete(users, *socketClient)
		log.Println("Number of client still connected ...", len(users))
		return
	}

	for client := range users {
		if err = client.websocket.WriteMessage(messageType, message); err != nil {
			log.Println("Cloud not send Message to ", client.ClientIP, err.Error())
		}
	}
}

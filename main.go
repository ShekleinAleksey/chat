package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chat/pkg/handler"
	"chat/pkg/repository"

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

type Message struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func newConnectUser(ws *websocket.Conn, ClientIP string) *ConnectUser {
	return &ConnectUser{
		websocket: ws,
		ClientIP:  ClientIP,
	}
}

func main() {

	log.Println("its allright")
	// repos := repository.NewRepository(db)
	// services := service.NewService(repos)
	// handlers := handler.NewHandler(services)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func init() {
	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, PATCH, DELETE, GET")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		IndexHandler(w, r)
	})

	http.HandleFunc("/ws", WebSocketHandler)
	http.HandleFunc("/sign-in", handler.SignIn)
	http.HandleFunc("/sign-up", handler.SignUp)
	http.HandleFunc("/getMessages", handler.GetMessages)
	http.HandleFunc("/getMessagesByUser", handler.GetMessages)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.Upgrade(w, r, nil)
}

var users = make(map[ConnectUser]int)

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err.Error())
		return
	}

	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	if ws == nil {
		log.Println("Websocket is nil")
		return
	}

	log.Println("Client connected:", ws.RemoteAddr().String())

	var socketClient *ConnectUser = newConnectUser(ws, ws.RemoteAddr().String())
	users[*socketClient] = 0
	log.Println("Number client connected ...", len(users))

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Ws disconnect waiting", err.Error())
			delete(users, *socketClient)
			log.Println("Number of client still connected ...", len(users))
			return
		}
		log.Println("message", string(message))
		t := time.Now().Format("2006-01-02 15:04:05")
		messageTemp := Message{Timestamp: t, Message: string(message)}
		jsonMessage, err := json.Marshal(messageTemp)
		if err != nil {
			log.Println(err)
			return
		}

		db, err := repository.NewDatabase()
		if err != nil {
			log.Fatalf("could not initialize database connection: %s", err)
		}
		chat_id := 1
		user_id := 5
		query := "INSERT INTO message_table (text, timestamp, chat_id, user_id) values ($1, $2, $3, $4) RETURNING id"
		row := db.QueryRow(query, string(message), t, chat_id, user_id)
		var id int
		if err := row.Scan(&id); err != nil {
			log.Fatal(err)
		}
		// timestamp := []byte(t)
		// ws.WriteMessage(websocket.TextMessage)
		// message = bytes.Join([][]byte{timestamp, message}, nil)

		for client := range users {
			if err = client.websocket.WriteMessage(messageType, jsonMessage); err != nil {
				log.Println("Cloud not send Message to ", client.ClientIP, err.Error())
			}
		}
	}

}

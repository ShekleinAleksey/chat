package handler

import (
	"chat/db"
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      string `json:"text"`
	ChatID    int    `json:"chat_id" db:"chat_id"`
	UserID    int    `json:"user_id" db:"user_id"`
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	log.Println("start")
	var messages []Message
	db, err := db.NewDatabase()
	if err != nil {
		log.Printf("could not initialize database connection: %s", err)
	}
	defer db.Close()
	query := "SELECT * FROM message_table"
	log.Println("2")
	err = db.Select(&messages, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("messages", messages)

	json.NewEncoder(w).Encode(messages)
}

func GetMessagesByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var messages []Message
	db, err := db.NewDatabase()
	if err != nil {
		log.Printf("could not initialize database connection: %s", err)
	}
	defer db.Close()
	query := "SELECT * FROM message_table m JOIN user_table u WHERE m.user_id = u.id"
	err = db.Select(&messages, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}

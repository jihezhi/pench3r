package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ChatMessage struct {
	Name string
	Content string
}

func randName(n int) string {
	name := make([]byte, n)
	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}
	return string(name)
}

func main() {
	requestSender := make(chan ChatMessage)
	var globalResponseMap map[string]*websocket.Conn
	globalResponseMap = make(map[string]*websocket.Conn)
	var historyChatContent []ChatMessage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	})
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		username := randName(6)
		fmt.Println(username)
		cm := ChatMessage{username, username + " enter."}
		globalResponseMap[username] = conn
		requestSender <- cm
		for _,m := range historyChatContent {
			conn.WriteJSON(m)
		}
		for {
			err := conn.ReadJSON(&cm)
			if err != nil {
				delete(globalResponseMap, cm.Name)
				cm.Content = username + " exit."
				requestSender <- cm
				return
			}
			requestSender <- cm
			// conn.WriteJSON(cm)
		}
	})
	go func() {
		for r := range requestSender {
			historyChatContent = append(historyChatContent, r)
			for _,c := range globalResponseMap {
				c.WriteJSON(r)
			}
		}
		close(requestSender)
	}()
	http.ListenAndServe(":8888", nil)
}

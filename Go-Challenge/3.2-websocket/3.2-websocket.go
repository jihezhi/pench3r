package main

// my.cnf
// character_set_server = utf8

// create database chatroom
// grant all privileges on chatroom.* to 'chat_admin'@'%' identified by 'chat_admin';
// flush privileges;
// use chatroom;
// create table chat_content(id int not null primary key auto_increment, uname varchar(10), content varchar(100))charset=utf8;

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

var db_conn *sql.DB

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ChatMessage struct {
	Name string
	Content string
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func mysqlInit() {
    db_conn, _ = sql.Open("mysql", "chat_admin:chat_admin@tcp(192.168.137.128:3306)/chatroom?charset=utf8")
    db_conn.SetMaxOpenConns(2000)
    db_conn.SetMaxIdleConns(1000)
    db_conn.Ping()
}

func getChatMsgByMysql() []map[string]string {
    rows, err := db_conn.Query("select uname, content from chat_content");
    checkErr(err)
    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for i := range scanArgs {
        scanArgs[i] = &values[i]
    }
    var record []map[string]string
    for rows.Next() {
        rows.Scan(scanArgs...)
        row := make(map[string]string)	
        for i, col := range values {
            row[columns[i]] = string(col.([]byte))
        }
        record = append(record, row)
    }
    return record
}

func insertChatMsgToMysql(uname, content string) {
	stmt, err := db_conn.Prepare("insert into chat_content(uname, content) values(?,?)")
	checkErr(err)
	_, err = stmt.Exec(uname, content)
	checkErr(err)
	stmt.Close()
}

func convertMapToChatMessage(map_ss []map[string]string) []ChatMessage {
	var chatmessage_slice []ChatMessage
	for _,s := range map_ss {
		cm := ChatMessage{s["uname"], s["content"]}
		chatmessage_slice = append(chatmessage_slice, cm)
	}
	return chatmessage_slice
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
	// mysql init
	mysqlInit()
	var historyChatContent []ChatMessage
	historyChatContent = convertMapToChatMessage(getChatMsgByMysql())
	fmt.Println(historyChatContent)
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
			go insertChatMsgToMysql(r.Name, r.Content)
			historyChatContent = append(historyChatContent, r)
			for _,c := range globalResponseMap {
				c.WriteJSON(r)
			}
		}
		close(requestSender)
	}()
	http.ListenAndServe(":8888", nil)
}

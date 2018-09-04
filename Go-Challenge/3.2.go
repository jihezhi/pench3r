// 解决了几个关键的问题，现在可以正常跑，但是写起来确实比较麻烦
// 接着再使用websocket来实现一下
package main

import (
    "fmt"
	"math/rand"
    "net"
    "log"
    "io"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ChatMessage struct {
	username string
	content string
	conn_state bool
}

func randName(n int) string {
	name := make([]byte, n)
	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}
	return string(name)
}

func handleConn(c net.Conn, historyContent []string) (string, chan ChatMessage, chan string) {
    requestChannel := make(chan ChatMessage)
    responseChannel := make(chan string)
	username := randName(6)
    tmp := make([]byte, 256)
	message := ChatMessage{username, username + " enter\n", true}
	// first write history content
	for _,msg := range historyContent {
		c.Write([]byte(msg))
	}
	// for loop to recive user input
    go func() {
		requestChannel <- message
        for {
            n, err := c.Read(tmp)
            if err != nil {
                if err != io.EOF {
                    fmt.Println("read error", err)
                    break
                }
            }
			// if send msg len is 0, exit and close conn
			if n == 0 {
				break
			}
			message.content = username + ": " + string(tmp[:n])
            requestChannel <- message
        }
		message.conn_state = false
		message.content = message.username + " exited\n"
		requestChannel <- message
        c.Close()
    }()
	// for loop to write response
    go func() {
        for {
			if !message.conn_state {
				break
			}
            if v, ok := <-responseChannel; ok {
                _, err := c.Write([]byte(v))
				if err != nil {
					break
				}
            } else {
				break
            }
        }
		close(responseChannel)
    }()
    return username, requestChannel, responseChannel
}

func main() {
    fmt.Println("server start: ")
    ln, err := net.Listen("tcp", ":8888")
    if err != nil {
        log.Fatal("net listen: ", err)
    }
    globalMessageChannel := make(chan ChatMessage)
	var storeRespChan map[string]chan string
	storeRespChan = make(map[string]chan string)
	var historyContent []string
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("net accept: ", err)
        }
        uname,req,res := handleConn(conn, historyContent)
		storeRespChan[uname] = res
        // fan-in
        go func(req chan ChatMessage) {
            for r := range req {
                globalMessageChannel <- r
                fmt.Printf("Recv msg: %v", r)
            }
        }(req)
        // fan-out
        go func() {
            for g := range globalMessageChannel {
				historyContent = append(historyContent, g.content)
                fmt.Printf("broadcast msg: %v", g.content)
				if !g.conn_state {
					delete(storeRespChan, g.username)
				}
                for _,c := range storeRespChan {
					c <- g.content
                }
            }
        }()
    }
}

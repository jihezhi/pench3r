// 整体的程序框架基本已经定型

package main

import (
    "fmt"
	"math/rand"
    "net"
    "log"
    "io"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 用于随机生成用户姓名
func randName(n int) string {
	name := make([]byte, n)
	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}
	return string(name)
}

func handleConn(c net.Conn, historyContent []string) (chan string, chan string) {
    requestChannel := make(chan string)
    responseChannel := make(chan string)
	username := randName(6)
    tmp := make([]byte, 256)
	// first write history content
	for _,msg := range historyContent {
		c.Write([]byte(msg))
	}
	// for loop to recive user input
    go func() {
        for {
            n, err := c.Read(tmp)
            if err != nil {
                if err != io.EOF {
                    fmt.Println("read error", err)
                    c.Close()
                    break
                }
            }
            requestChannel <- username + ": " + string(tmp[:n])
        }
        c.Close()
    }()
	// for loop to write response
    go func() {
        for {
            if v, ok := <-responseChannel; ok {
                c.Write([]byte(v))
            } else {
                return
            }
        }
    }()
    return requestChannel, responseChannel
}

func main() {
    fmt.Println("server start: ")
    ln, err := net.Listen("tcp", ":8888")
    if err != nil {
        log.Fatal("net listen: ", err)
    }
    globalMessageChannel := make(chan string)
    // 用于保存每个了解的responsechannel
	var storeResponseChannel []chan string
    // 用于保存所有用户发送过的信息
	var historyContent []string
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("net accept: ", err)
        }
        req,res := handleConn(conn, historyContent)
        storeResponseChannel = append(storeResponseChannel, res)
        // fan-in
        go func(req chan string) {
            for r := range req {
                globalMessageChannel <- r
                fmt.Printf("Recv msg: %v", r)
            }
        }(req)
        // fan-out
        go func() {
            for g := range globalMessageChannel {
				historyContent = append(historyContent, g)
                fmt.Printf("broadcast msg: %v", g)
                for _,c := range storeResponseChannel {
                    c <- g
                }
            }
        }()
    }
}

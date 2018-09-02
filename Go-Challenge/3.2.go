// 整体的程序框架基本已经定型
package main

import (
    "fmt"
    "net"
    "log"
    "io"
)

func handleConn(c net.Conn) (chan string, chan string) {
    requestChannel := make(chan string)
    responseChannel := make(chan string)
    tmp := make([]byte, 256)
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
            requestChannel <- c.RemoteAddr().String() + string(tmp)
        }
        c.Close()
    }()
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
    storeResponseChannel := make([]chan string, 0)
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("net accept: ", err)
        }
        req,res := handleConn(conn)
        storeResponseChannel = append(storeResponseChannel, res)
        go func(req chan string) {
            for r := range req {
                globalMessageChannel <- r
                fmt.Printf("Recv msg: %v", r)
            }
        }(req)
        go func() {
            for g := range globalMessageChannel {
                fmt.Printf("broadcast msg: %v", g)
                for _,c := range storeResponseChannel {
                    c <- g
                }
            }
        }()
    }
}

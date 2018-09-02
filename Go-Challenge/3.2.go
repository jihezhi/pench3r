// 未完成的一个简单demo,需要再调试一下
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
    go func() {
        tmp := make([]byte, 256)
        for {
            _, err := c.Read(tmp)
            if err != nil {
                if err != io.EOF {
                    fmt.Println("read error", err)
                }
                break
            }
            requestChannel <- string(tmp)
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
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("net accept: ", err)
        }
        req,res := handleConn(conn)
        go func(req chan string) {
            for r := range req {
                globalMessageChannel <- r
                fmt.Printf("Recv msg: %v", r)
            }
        }(req)
        go func(res chan string) {
            for g := range globalMessageChannel {
                res <- g
                fmt.Printf("broadcast msg: %v", g)
            }
        }(res)
    }
}

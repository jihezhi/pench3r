// 跳跃了一下。
package main

import (
    "net/http"
    "log"
)

func main() {
    e := http.ListenAndServe(":8888", http.FileServer(http.Dir(".")))
    if e != nil {
        log.Fatal("ListenAndServe: ",e)
    }
}

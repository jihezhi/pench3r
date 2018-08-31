// 跳跃了一下。
package main

import (
    "net/http"
    "log"
)

func main() {
    log.Fatal(http.ListenAndServe(":8888", http.FileServer(http.Dir("."))))
}

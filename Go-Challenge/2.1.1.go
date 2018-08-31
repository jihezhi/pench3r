// 重新更新一下，不舍得删旧代码。
package main

import (
    "fmt"
    "net/http"
    "strconv"
)

func main() {
    requestSender := make(chan int)
    responseReceiver := make(chan int)
    sumChannel := make(chan int)
    var sum int
    // 初始化源管道中的数据
    go func() { sumChannel <- 0}()
    go func() {
        for s := range sumChannel {
            sum = s + <-requestSender
            responseReceiver <- sum
            go func(s int) {sumChannel <- s}(sum)
        }
    } ()

    http.HandleFunc("/getValue", func(w http.ResponseWriter, r *http.Request) {
        // 获取querystring，并将其转化为int类型
        s := r.URL.Query()["int"][0]
        sum, _ := strconv.Atoi(s)
        requestSender <- sum
        fmt.Fprint(w, <-responseReceiver)
    })
    http.ListenAndServe(":8888", nil)
}

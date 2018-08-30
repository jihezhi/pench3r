// 第一版，写完了注释还是感觉有点虚，不知道理解对不对
package main

import (
    "fmt"
    "net/http"
    "strconv"
)

var src_chan chan int
var dst_chan chan int

// 函数功能：将from channel中的数据于s相加并存入to channel
func add_via_channel(to chan int, s int, from chan int) {
    to <- s + <- from
}

func main() {
    src_chan := make(chan int)
    dst_chan := make(chan int)
    // 初始化源管道中的数据
    go func() { src_chan <- 0}()
    http.HandleFunc("/getValue", func(w http.ResponseWriter, r *http.Request) {
        // 获取querystring，并将其转化为int类型
        s := r.URL.Query()["int"][0]
        sum, _ := strconv.Atoi(s)
        // 将相加后的结果放入dst_chan中，在高并发中只有当src可用时才会执行
        // 因此会通过src的阻塞进行并发的控制，这里就相当于mutex.lock
        go add_via_channel(dst_chan, sum, src_chan)
        // 后面的操作相当于原子操作不会被其他所影响
        // 提取出结果并输出
        sum = <-dst_chan
        fmt.Fprint(w, sum)
        // 再将结果放入到src中，此时dst channel中为空，这里就相当于mutex.unlock
        go func(s int) { src_chan <- s }(sum)
    })
    http.ListenAndServe(":8888", nil)
}

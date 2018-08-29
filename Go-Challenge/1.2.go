// 感觉有点跑题了
package main

import (
    "fmt"
    "sort"
)

func main() {
    t := []int{2, 4, 1, 3, 7, 9}
    sort.Ints(t)
    fmt.Println(t)
}

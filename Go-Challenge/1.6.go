package main

import (
    "sort"
    "fmt"
)

func main() {
    float_map := make(map[float64]string)
    float_map[3.14] = "a"
    float_map[2.76] = "b"
    float_map[7.89] = "c"
    float_map[5.34] = "d"
    fmt.Println(float_map)    // map[2.76:b 7.89:c 5.34:d 3.14:a] 
    var sorted_key []float64
    for k := range float_map {
        sorted_key = append(sorted_key,k)
    }
    sort.Float64s(sorted_key)
    var sorted_value []string
    for _,k := range sorted_key {
        sorted_value = append(sorted_value, float_map[k])
    }
    fmt.Println(sorted_value)   // [b a d c]
}

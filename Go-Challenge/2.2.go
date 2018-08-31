package main

import (
    "fmt"
    "math/rand"
)

type determinant [3][3]int

func generate_dt() determinant {
    return determinant{[3]int{rand.Intn(10),rand.Intn(10),rand.Intn(10)},[3]int{rand.Intn(10),rand.Intn(10),rand.Intn(10)},[3]int{rand.Intn(10),rand.Intn(10),rand.Intn(10)}}
}

func determin_sum(dt determinant) int {
    return 2
}

func generate_subtask(dt_array []determinant, size int) chan int {
    c := make(chan int)
    go func () {
        for i := 0; i<size; i++ {
            c <- determin_sum(dt_array[i])
        }
    }()
    return c
}

func fanin_pattern(input1, input2, input3, input4 chan int) chan int {
    c := make(chan int)
    go func() {
        for {
            select {
                case i := <-input1: c <- i
                case i := <-input2: c <- i
                case i := <-input3: c <- i
                case i := <-input4: c <- i
            }
        }
    }()
    return c
}

func main() {
    var dt_array [100]determinant
    for i := 0; i<100; i++ {
        dt_array[i] = generate_dt()
    }

    result := fanin_pattern(generate_subtask(dt_array[:25], 25), generate_subtask(dt_array[25:50], 25), generate_subtask(dt_array[50:75], 25), generate_subtask(dt_array[75:100], 25))

    var sum = 0
    for i := 0; i<100; i++ {
        sum = sum + <-result
    }

    fmt.Println(sum)
}            

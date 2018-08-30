package main

import (
    "io"
    "os"
)

func main() {
    f,_ := os.Open("file.txt")
    io.Copy(os.Stdout, f)
}

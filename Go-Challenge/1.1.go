package main

import (
    "io"
    "os"
    "strconv"
    "bytes"
    "errors"
)

func main() {
    // string
    s := "first"
    io.WriteString(os.Stdout, s)
    
    // int
    i := 123
    bin := []byte(strconv.Itoa(i))
    io.Copy(os.Stdout, bytes.NewReader(bin))
    
    // float64
    var f float64 = 1.2345
    io.WriteString(os.Stdout, strconv.FormatFloat(f, 'f', 4, 64))
    
    // [3]bool
    ba := [3]bool{true, false, true}
    for _,v := range ba {
        io.WriteString(os.Stdout, strconv.FormatBool(v))
    }
    
    // error
    e := errors.New("error")
    io.WriteString(os.Stdout, e.Error())
}

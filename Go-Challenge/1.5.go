// 因为之前看到x7night使用multireader了，所以只能另换写法了...

package main

import (
    "io"
    "os"
    "bytes"
)

func main() {
    // string
    s := "string"
    io.WriteString(os.Stdout, s)

    // []byte
    var buf bytes.Buffer
    buf.Write([]byte("[]byte"))
    buf.WriteTo(os.Stdout)

    // file
    f, _ := os.Open("file.txt")
    io.Copy(os.Stdout, f)
}

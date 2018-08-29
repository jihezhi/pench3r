package main

import (
     "fmt"
	   "strconv"
	   "encoding/hex"
	   )

func main() {
  // 1
	i := 12345
	// 10 to 16
	h := strconv.FormatInt(int64(i), 16)
  fmt.Println(h)
  
	// 16 to byte
  var dest []byte
	fmt.Sscanf(h, "%X", &dest)
	fmt.Println(string(dest))
	
	// second method
	s, _ := hex.DecodeString(h)
	fmt.Printf("%s",s)
    
  // 2
	s1 := "3.1415926"
	f,_ := strconv.ParseFloat(s1, 64)
	fmt.Printf("%T value is %v",f, f)

}

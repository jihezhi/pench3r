// 粗糙的一版
package main

import (
    "fmt"
)

type eval_desc_I interface {
    evaluate() int
    describe() string
}

type myType struct {
    age int
    name string
}

func (mt myType) evaluate() int {
    return mt.age
}

func (mt myType) describe() string {
    return mt.name
}

func multitype_show(typename ...eval_desc_I) {
    if len(typename) == 0 {
        fmt.Println("you need pass one eval_desc_I interface at least")
        return
    }
    for _,v := range typename {
        fmt.Println(v.evaluate(), v.describe())
    }
}

func main() {
    mt := myType{18, "pench3r"}
    mt1 := myType{19, "pench3r"}
    mt2 := myType{20, "pench3r"}
    multitype_show(mt, mt1, mt2)
    multitype_show()
}

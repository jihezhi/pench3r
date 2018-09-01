// 粗糙的一版
package main

import (
    "fmt"
    "errors"
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

func multitype_show(typename ...interface{}) error {
    if len(typename) == 0 {
        return errors.New("you need pass one eval_desc_I interface at least")
    }
    for _,v := range typename {
        switch t := v.(type) {
            case myType:
                fmt.Println(t, v.(myType).evaluate(), v.(myType).describe())
            default:
                return fmt.Errorf("you pass type is %T,need pass type of implement eval_desc_I interface", t)
        }
    }
    return nil
}


func main() {
    mt := myType{18, "pench3r"}
    mt1 := myType{19, "pench3r"}
    mt2 := myType{20, "pench3r"}
    multitype_show(mt, mt1, mt2)
    err := multitype_show(1)
    if err != nil {
        fmt.Println(err)
    }
}

package test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestDemo(t *testing.T) {
	fmt.Println(strings.Join([]string{"1", "@", "2"}, ","))
	fmt.Println(strconv.FormatInt(12, 10))
	fmt.Println(rune(12))
	//new(T).Hello()
	//new(C2).Hello()
	//var c Class
	//c = new(C2)
	//c.Hello()
}

//
//type Class interface {
//	Hello()
//}
//
//type T int
//
//func (t *T)Hello()  {
//	fmt.Println(t)
//}
//
//type C2 struct {
//	T
//}
//func (c *C2)Hello()  {
//
//	fmt.Println(c)
//}

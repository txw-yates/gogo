package main

import (
	"fmt"
	"reflect"
)

type gg struct {
	name string
}

func (g *gg) Init() {
	fmt.Println("init")
}

type ci interface {
	Init()
}

func main() {
	hh := &gg{
		name: "ggg",
	}

	c := reflect.Indirect(reflect.ValueOf(hh)).Type()
	fmt.Println(c.Name())

	g := reflect.New(c)

	//g := reflect.ValueOf(hh)
	g.MethodByName("Init").Call(nil)
}
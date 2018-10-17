package main

import (
		"net/http"
		"log"
		"fmt"
)

type TestController struct {
	BaseController
	name string
}

func (test *TestController) Get(name string) {
	fmt.Println(test.Context.Request.URL.Query())
	test.Context.Response.Write([]byte(name))
}

func main() {
	rg := new(RouterRegister)

	rg.Add("/user/:name", TestController{})

	err := http.ListenAndServe(":9090", rg)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
package main

import "net/http"

// 控制器的基类
type ControllerInterface interface {}

// 参数
type Context struct {
	Request *http.Request
	Response http.ResponseWriter
}

type BaseController struct {
	Context *Context
}
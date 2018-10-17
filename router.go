package main

import (
	"regexp"
	"reflect"
	"strings"
	"net/http"
	"net/url"
)

// 路由信息
type Router struct {
	Regex *regexp.Regexp
	Params map[int]string
	Controller reflect.Type
}

// 路由注册
type RouterRegister struct {
	Routers []*Router
}

func (routerRegister *RouterRegister) Add(pattern string, c ControllerInterface) {
	parts := strings.Split(pattern, "/")
	params := make(map[int]string)

	j := 0
	for i, part := range parts {
		// 路由/user/:id([0-9]+)
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			// 判断是否对参数有限制
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			// 去除":"
			params[j] = part[1:]
			parts[i] = expr
			j++
		}
	}

	// 重新join，生成新的路由规则/user/([0-9]+)
	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
		return
	}

	router := &Router{
		Params: params,
		Regex: regex,
		// 包名.struct名称, 如main.struct
		Controller: reflect.Indirect(reflect.ValueOf(c)).Type(),
	}

	routerRegister.Routers = append(routerRegister.Routers, router)
}

func (routerRegister *RouterRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 是否有匹配
	match := false

	requestURL := r.URL.Path
	// 循环routers
	for _, router := range routerRegister.Routers {
		// 判断正则是否能匹配上
		if !router.Regex.MatchString(requestURL) {
			continue
		}
		matchs := router.Regex.FindStringSubmatch(requestURL)
		// 如果匹配的结果不与整个请求url相同，跳过
		if matchs[0] != requestURL {
			continue
		}

		// 更改匹配标志
		match = true

		// 将匹配的值与对应的参数匹配
		params := make(map[string]string)
		if len(router.Params) > 0 {
			// 取出url中query值
			values := r.URL.Query()
			// 将匹配到的参数加入其中
			for i, val := range matchs[1:] {
				values.Add(router.Params[i], val)
				params[router.Params[i]] = val
			}

			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
		}

		// 根据反射创建新的struct
		c := reflect.New(router.Controller)

		// context赋值
		context := &Context{r, w}
		c.Elem().FieldByName("Context").Set(reflect.ValueOf(context))

		in := make([]reflect.Value, 1)

		in[0] = reflect.ValueOf(params["name"])
		if r.Method == "GET" {
			method := c.MethodByName("Get")
			method.Call(in)
		}
	}

	// 没有匹配
	if !match {
		http.NotFound(w, r)
	}
}
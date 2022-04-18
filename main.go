package coredemo

import (
	"coredemo/framework"
	"net/http"
)

func main() {
<<<<<<< HEAD
	server := &http.Server{//这个结构体的Handler字段就只去使用自己的Handler函数
		// 自定义的请求核心处理函数
		Handler: framework.NewCore(),//实际上server结构体的的一个方法接口就是Handler,但是Handler只有一个方法
		//就是serverHTTP,只要我返回一个结构体,这个结构体有Handler这个方法,就符合这个接口
		// 请求监听地址
		Addr:    ":8080",
=======
	server := &http.Server{ //这个结构体的Handler字段就只去使用自己的Handler函数
		// 自定义的请求核心处理函数
		Handler: framework.NewCore(), //实际上server结构体的的一个方法接口就是Handler,但是Handler只有一个方法
		//就是serverHTTP,只要我返回一个结构体,这个结构体有Handler这个方法,就符合这个接口
		// 请求监听地址
		Addr: ":8080",
>>>>>>> 8e475c8 (Initial commit)
	}
	server.ListenAndServe()
}

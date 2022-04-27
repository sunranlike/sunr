package framework

/*
import (
	"encoding/json"
	"html/template"
)

// IResponse代表返回方法
type IResponse interface {
	// Json输出
	IJson(obj interface{}) IResponse

	// Jsonp输出
	IJsonp(obj interface{}) IResponse

	//xml输出
	IXml(obj interface{}) IResponse

	// html输出
	IHtml(template string, obj interface{}) IResponse

	// string
	IText(format string, values ...interface{}) IResponse

	// 重定向
	IRedirect(path string) IResponse

	// header
	ISetHeader(key string, val string) IResponse

	// Cookie
	ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	// 设置状态码
	ISetStatus(code int) IResponse

	// 设置200状态
	ISetOkStatus() IResponse
}
func (ctx *Context) IJsonp(obj interface{}) IResponse {
	// 获取请求参数 callback
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.ISetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成 XSS 攻击
	//XSS攻击通常指的是通过利用网页开发时留下的漏洞，通过巧妙的方法注入恶意指令代码到网页，
	//使用户加载并执行攻击者恶意制造的网页程序。这些恶意网页程序通常是JavaScript，
	//但实际上也可以包括Java，VBScript，
	//ActiveX，Flash或者甚至是普通的HTML。
	//攻击成功后，攻击者可能得到更高的权限（如执行一些操作）、私密网页内容、
	//会话和cookie等各种内容。
	callback := template.JSEscapeString(callbackFunc)

	// 输出函数名
	_, err := ctx.responseWriter.Write([]byte(callback))
	if err != nil {
		return ctx
	}
	// 输出左括号
	_, err = ctx.responseWriter.Write([]byte("("))
	if err != nil {
		return ctx
	}
	// 数据函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.responseWriter.Write(ret)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.responseWriter.Write([]byte(")"))
	if err != nil {
		return ctx
	}
	return ctx
}

func (ctx *Context) IHtml(file string, obj interface{}) IResponse {
	// 第一步,读取模版文件，创建 template 实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 第二步,执行 Execute 方法将 obj 和模版进行结合
	//两步结合就可以实现模板+数据的动态模型
	if err := t.Execute(ctx.responseWriter, obj); err != nil {
		return ctx
	}

	ctx.ISetHeader("Content-Type", "application/html")
	return ctx
}

*/

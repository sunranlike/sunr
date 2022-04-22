package framework

import "mime/multipart"
import "github.com/spf13/cast"

const defaultMultipartMemory = 32 << 20 // 32 MB
// 代表请求包含的方法
type IRequest interface {
	// QueryInt 请求地址 url 中带的参数
	// 形如: foo.com?a=1&b=bar&c[]=bar
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	// 路由匹配中带的参数
	// 形如 /book/:id
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	// form 表单中带的参数
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	// json body
	BindJson(obj interface{}) error

	// xml body
	BindXml(obj interface{}) error

	// 其他格式
	GetRawData() ([]byte, error)

	// 基础信息
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	// header
	Headers() map[string][]string
	Header(key string) (string, bool)

	// cookie
	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

// 获取请求地址中所有参数
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

// 获取 Int 类型的请求参数
func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			// 使用 cast 库将 string 转换为 Int
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}
func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) QueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0], true
		}
	}
	return def, false
}

func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (ctx *Context) Query(key string) interface{} {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals[0]
	}
	return nil
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		ctx.request.ParseForm()
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}
	return def, false
}

func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil {
		if err := ctx.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := ctx.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, err
}

func (ctx *Context) Form(key string) interface{} {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals[0]
		}
	}
	return nil
}

// 路由匹配中带的参数
// 形如 /book/:id
func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if val := ctx.Param(key); val != nil {
		// 通过cast进行类型转换
		return cast.ToInt(val), true
	}
	return def, false
}

func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt64(val), true
	}
	return def, false
}

func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToFloat64(val), true
	}
	return def, false
}

func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToFloat32(val), true
	}
	return def, false
}

func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToBool(val), true
	}
	return def, false
}

func (ctx *Context) ParamString(key string, def string) (string, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToString(val), true
	}
	return def, false
}

// 获取路由参数
func (ctx *Context) Param(key string) interface{} {
	if ctx.params != nil {
		if val, ok := ctx.params[key]; ok {
			return val
		}
	}
	return nil
}

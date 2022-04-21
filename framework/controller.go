package framework

// ControllerHandler ControllerHandler其实就是个函数的别名
//它被用来1存入一个slice中,方便链式调用
type ControllerHandler func(c *Context) error

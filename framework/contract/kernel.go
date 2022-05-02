package contract

import "net/http"

// KernelKey 提供 kenel 服务凭证
const KernelKey = "hade:kernel"

// Kernel 接口提供框架最核心的结构，kernel是内核的意思，确实本接口是实现web服务，自然是我们的核心
type Kernel interface {
	// HttpEngine http.Handler结构，作为net/http框架使用, 实际上是gin.Engine
	//返回值也是个接口，是http.Handler，这个接口只有一个方法就是ServeHttp，所以我们返回一个
	//GIN.ENGINE也是可以的，因为ENGINE实现了ServeHttp方法，也是实现了http.Handler方法
	//那这里为什么使用的返回值是个http.Handler?因为想要更加抽象,不要单独以来gin框架,后期可以修改框架
	HttpEngine() http.Handler
}

package demo

import (
	demoService "github.com/sunranlike/hade/app/provider/demo"
	"github.com/sunranlike/hade/framework/gin"
)

type DemoApi struct {
	service *Service
}

//这里将路由注册,当然这里只是个demo
func Register(r *gin.Engine) error {
	api := NewDemoApi() //这是个啥结构?这是注册路由的规范吗?这个原来是个存放handler的地方
	//engine当然有Bind方法,因为engine也实现了framework.Container接口
	//我们在main方法中绑定的是框架服务,这里绑定的是业务服务,当然这些个业务服务都要实现serviceProvider接口
	r.Bind(&demoService.DemoProvider{})

	r.GET("/demo/demo", api.Demo) //绑定处理器
	r.GET("/demo/demo2", api.Demo2)
	r.POST("/demo/demo_post", api.DemoPost)
	return nil
}

func NewDemoApi() *DemoApi { //返回一个DemoApi结构,这是个啥?
	service := NewService() //NewService是个啥?又绑定到了Service中了?
	return &DemoApi{service: service}
}

// Demo godoc
// @Summary 获取所有用户
// @Description 获取所有用户
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo [get]
func (api *DemoApi) Demo(c *gin.Context) {
	//appService := c.MustMake(contract.AppKey).(contract.App)
	//baseFolder := appService.BaseFolder()
	users := api.service.GetUsers()
	usersDTO := UserModelsToUserDTOs(users)
	c.JSON(200, usersDTO)
}

// Demo godoc
// @Summary 获取所有学生
// @Description 获取所有学生
// @Produce  json
// @Tags demo
// @Success 200 array []UserDTO
// @Router /demo/demo2 [get]
func (api *DemoApi) Demo2(c *gin.Context) {
	demoProvider := c.MustMake(demoService.DemoKey).(demoService.IService)
	students := demoProvider.GetAllStudent()
	usersDTO := StudentsToUserDTOs(students)
	c.JSON(200, usersDTO)
}

func (api *DemoApi) DemoPost(c *gin.Context) {
	type Foo struct {
		Name string
	}
	foo := &Foo{}
	err := c.BindJSON(&foo)
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, nil)
}

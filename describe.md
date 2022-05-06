##app文件夹
app文件夹干了什么？他存放了业务级别服务应用
##config文件夹
存放配置文件
##framework文件夹
+ 存放主题框架，融入的包，比如gin，cobra
+ 实现业务级别服务，contract+provider
+ 融入gin和cobra并且融入container
+ 实现中间件middleware
+ 一些工具包
+ container文件和provider文件定义了容器接口和服务提供者接口，并且定义了一个实现了container接口的实例
+ util实现了一些工具
##provider文件夹
业务级别

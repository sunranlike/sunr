本文件夹对框架服务做了实现，注意是框架服务，都是基础服务。
一个contract文件对应一个provider和service
#provider
providerr文件干了什么：定义xxxProvider结构体，**实现ServiceProvider接口**
+ 定义xxxProvider结构体，并以他为基实现service provider接口的五个方法。
+ 没了=。=

#service
service文件干什么：**服务真正的实现**
+ 定义xx结构体，围绕他实现contract，实现在contract协议件的具体实现

#总结
+ 这样就把框架级别服务切分开了三个部分:**服务接口定义,实现Provider接口,实现服务接口**
+ 在framework.contract文件夹下实现了协议描绘，定义了key和和具体服务实现接口
+ 在framework.provider文件夹中为每个框架级别服务新建文件夹,在这个文件夹中有两个文件,provider和service
+ provider文件负责**实现serviceProvider接口**(5个方法)
+ service文件中负责实现**具体服务实现接口**(在contract中定义的)
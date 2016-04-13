### GinWeb

A golang web project template based on gin framework

### Features:

gin本身已经相当灵活,但是由于项目中有一些模块是我们必须或者经常使用的,所以进行了集成,如果不需要这些模块或者要使用自己的模块实现,请考虑直接使用gin框架

1. 集成[endless](https://github.com/fvbock/endless)在类Unix平台下实现不停机重启服务
2. 集成log4go: 类似于log4j的配置式的log输出,可以控制日记级别,输出格式以及滚动
3. 集成recovery: 出错时尝试进行恢复并使用log4go记录为critical级别

### USAGE:

```go
func main() {
    r := ginweb.New()
    
    router.GET("/welcome", func(c *gin.Context) {
        firstname := c.DefaultQuery("firstname", "Guest")
        lastname := c.Query("lastname")
        c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
    })
    
    ginweb.Run(conf.Conf.SERVER.PORT, r)
}
```

后台启动app（供参考）：
```sh
$> nohup ./appname > /dev/null 2>stderr.log &
```

不停机重启服务（在在Unix/Linux/Mac）上能起作用：

```sh
$> pgrep appname | xargs kill -HUP
```

### TODO:

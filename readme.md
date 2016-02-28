### GinWeb

A golang web project template based on gin framework

### Features:

gin本身已经相当灵活,但是由于项目中有一些模块是我们必须或者经常使用的,所以进行了集成,如果不需要这些模块或者要使用自己的模块实现,请考虑直接使用gin框架

1. 集成log4go: 类似于log4j的配置式的log输出,可以控制日记级别,输出格式以及滚动
2. 集成recovery: 出错时尝试进行恢复并使用log4go记录为critical级别

### USAGE:

	func main() {
	    r := ginweb.New()
	    
	    router.GET("/welcome", func(c *gin.Context) {
            firstname := c.DefaultQuery("firstname", "Guest")
            lastname := c.Query("lastname")
            c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
        })
	    
	    ginweb.Run(conf.Conf.SERVER.PORT, r)
	}

### TODO:

1. mgo.v2
2. mgox.v2
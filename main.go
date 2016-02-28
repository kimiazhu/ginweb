// Description: GOTAKO的V2版本，此版本基于gin实现
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kimiazhu/ginweb/conf"
	"github.com/kimiazhu/ginweb/controller/apps"
	"github.com/kimiazhu/ginweb/server"
)

func main() {
	g := gin.New()

	g.Use(AccessLog())
	g.Use(Recovery())

	app := g.Group("/apps")
	{
		app.POST("checkupdate", apps.CheckUpdate)
	}

	server.Start(":"+Conf.SERVER.PORT, g)
}

// Description: GOTAKO的V2版本，此版本基于gin实现
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package main

import (
	"github.com/kimiazhu/ginweb/controller/apps"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/ginweb/server"
	. "github.com/kimiazhu/ginweb/conf"
)

func main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	app := r.Group("/apps")
	{
		app.GET("checkupdate", apps.CheckUpdate)
	}

	server.Start(":" + Conf.SERVER.PORT, r)

}

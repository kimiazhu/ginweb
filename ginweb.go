// Description: ginweb is light weight encapsulation of gin framework
// Author: ZHU HAIHUA
// Since: 2016-02-28 14:57
package ginweb

import (
	"github.com/gin-gonic/gin"
	. "github.com/kimiazhu/ginweb/midware"
	"github.com/kimiazhu/ginweb/server"
	"github.com/kimiazhu/ginweb/conf"
)

const VERSION = "0.0.1"

func New() *gin.Engine {
	gin.SetMode(conf.Conf.SERVER.RUNMODE)
	g := gin.New()


	g.Use(Recovery())
	g.Use(AccessLog())
	return g
}

func Run(port string, engin *gin.Engine) {
	server.Start(":"+port, engin)
}

// Description: ginweb is light weight encapsulation of gin framework
// Author: ZHU HAIHUA
// Since: 2016-02-28 14:57
package ginweb

import (
	"github.com/gin-gonic/gin"
	. "github.com/kimiazhu/ginweb/midware"
	"github.com/kimiazhu/ginweb/server"
	"github.com/kimiazhu/ginweb/conf"
	"github.com/kimiazhu/log4go"
	"sync"
)

const VERSION = "0.0.1"

// a component is different from middleware,
// it will be initialize just before the
// application running, then you can use it in
// the entire life-cycle of the application
type component struct {
	name string
	config interface{}
	initialize func(config interface{}) (error)
}

var initCompOnce sync.Once
var components []component

func New() *gin.Engine {
	gin.SetMode(conf.Conf.SERVER.RUNMODE)
	g := gin.New()
	g.Use(Recovery())
	g.Use(AccessLog())
	return g
}

func RegisterComponent(name string, config interface{}, initialize func(config interface{}) (error))  {
	components = append(components, component{name, config, initialize})
}

func Run(port string, engin *gin.Engine) {
	initCompOnce.Do(initialize)
	server.Start(":"+port, engin)
}

// initialize used to init all components before the app start
func initialize() {
	for _, c := range components {
		e := c.initialize(c.config)
		if e != nil {
			log4go.Error("initialize component [%s] error! %v", c.name, e)
		} else {
			log4go.Debug("initialize component [%s] success", c.name)
		}
	}
}
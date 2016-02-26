// Description: apps包提供跟App相关的服务
// version.go 提供跟App版本相关的服务，如：更新检查
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package apps

import (
	"github.com/gin-gonic/gin"
	"conf"
)

func CheckUpdate(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": conf.Conf,
	})
}
// Description: server 提供Windows下的服务器实现。
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	_server = &serverWin{}
}

type serverWin struct {
}

func (s *serverWin) Run(addr string, handler http.Handler) {
	handler.(*gin.Engine).Run(addr)
}

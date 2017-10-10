// Description: server 提供抽象的服务器接口。
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package server

import (
	"net/http"
)

type server interface {
	Run(addr string, handler http.Handler, args... interface{})
}

var _server server

// Start 用于启动服务器
func Start(addr string, handler http.Handler, args ...interface{}) {
	_server.Run(addr, handler, args...)
}

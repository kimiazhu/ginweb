// Description: server 提供Unix下的服务器实现。
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package server

import (
	"github.com/fvbock/endless"
	"net/http"
)

func init() {
	_server = &serverUnix{}
}

type serverUnix struct {
}

func (s *serverUnix) Run(addr string, handler http.Handler) {
	endless.ListenAndServe(addr, handler)
}

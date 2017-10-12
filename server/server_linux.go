// Description: server 提供Unix下的服务器实现。
// Author: ZHU HAIHUA
// Since: 2016-02-26 18:56
package server

import (
	"github.com/fvbock/endless"
	"net/http"
	"time"
)

func init() {
	_server = &serverLinux{}
}

type serverLinux struct {
}

func (s *serverLinux) Run(addr string, handler http.Handler, args ...interface{}) {
	if len(args) > 0 {
		endless.DefaultHammerTime = args[0].(time.Duration)
	}
	endless.ListenAndServe(addr, handler)
}

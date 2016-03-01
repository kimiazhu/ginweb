// Description: utils 服务器应答标准结构
// Author: ZHU HAIHUA
// Since: 2016-02-26 20:31
package model

import (
	"github.com/kimiazhu/ginweb/util"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Tag  string      `json:"tag,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (r *Response) Dump2Json() string {
	return util.Dump2Json(r)
}

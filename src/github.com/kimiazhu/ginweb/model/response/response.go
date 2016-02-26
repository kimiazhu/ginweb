// Description: utils 服务器应答标准结构
// Author: ZHU HAIHUA
// Since: 2016-02-26 20:31
package response

import (
	"encoding/json"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Tag  string      `json:"tag,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (r *Response) Dump2Json() string {
	return Dump2Json(r)
}

func Dump2Json(obj interface{}) string {
	result, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(result)
}

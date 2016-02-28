// Description: errcode 错误码定义
// Author: ZHU HAIHUA
// Since: 2016-02-26 20:37
package errcode

const (
	Success        = 0    //成功
	Unknown        = -1   //未知错误，一般性错误
	WrongArguments = -2   // 收到错误的参数
	Signature      = -200 // 签名错误
	Remote         = -100 // 远程系统调用一般性错误
)

// Description: utils 工具包 rand.go 提供随机数相关服务
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:08
package utils

import (
	"crypto/rand"
)

var charTable = []rune("abcdefghijkmnpqrstuvwxyz23456789")

// RandStrN 返回长度为N的随机字符和数字组合,其中不包含容易被混淆的[0/1/o/l]四个字符
func RandStrN(n int) string {
	random := make([]byte, n)
	result := make([]rune, n)
	rand.Read(random[:])
	for i := 0; i < len(random); i++ {
		result[i] = charTable[uint(random[i]>>3)]
	}
	return string(result)
}

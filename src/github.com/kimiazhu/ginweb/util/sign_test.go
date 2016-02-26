// Description: utils
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:08
package utils

import (
	"fmt"
	"testing"
)

func TestSignVerify2(t *testing.T) {
	jsStr := "{\"Key1\":\"Value1\",\"Key2\":\"\",\"Key3\":[\"Value3.1\",\"Value3.2\"],\"sign\":\"aef03ab309f230d32f612daee3fb882e4f533404\"}"
	r := SignVerify2(jsStr, Secret_Key_User_Center)
	fmt.Printf("\nresult is :%v\n", r)
}

func TestBuildSignedJsonStr(t *testing.T) {
	r, _ := BuildSignedJsonStr(map[string]interface{}{"Key1": "Value1", "Key2": "", "Key3": []string{"Value3.1", "Value3.2"}}, Secret_Key_User_Center)
	fmt.Printf("%s\n", r)
}

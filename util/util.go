// Description: util
// Author: ZHU HAIHUA
// Since: 2016-02-29 21:41
package util

import "encoding/json"

func Dump2Json(obj interface{}) string {
	result, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(result)
}

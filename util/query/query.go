// Description: query
// Author: ZHU HAIHUA
// Since: 2016-02-28 21:46
package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// MustInt try to extract value from request param by the givin key,
// if failed will raise a panic
func MustInt(c *gin.Context, key string) int {
	if strv := c.Query(key); strv != "" {
		v, e := strconv.Atoi(strv)
		if e != nil {
			panic(e)
		}
		return v
	} else {
		panic(fmt.Sprintf("cannot get value of param [%s]", key))
	}
}

// Int try to extract value from request param by the givin key,
// if failed the error will set and defaults will be returned (if exists)
func Int(c *gin.Context, key string, defaults ...int) (int, error) {
	if strv := c.Query(key); strv != "" {
		return strconv.Atoi(strv)
	} else if len(defaults) > 0 {
		return defaults[0], nil
	} else {
		return strconv.Atoi(strv)
	}
}

// Description: model
// Author: ZHU HAIHUA
// Since: 2016-02-28 20:59
package param

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// MustInt try to extract value from request param by the givin key,
// if failed will return the default value
func MustInt(c *gin.Context, key string, defaultValue int) int {
	strv := c.Param(key)
	v, e := strconv.Atoi(strv)
	if e != nil {
		v = defaultValue
	}
	return v
}

// Int try to extract value from request param by the givin key,
// if failed the error will set and defaults will be returned (if exists)
func Int(c *gin.Context, key string, defaults ...int) (int, error) {
	if strv := c.Param(key); strv != "" {
		return strconv.Atoi(strv)
	} else if len(defaults) > 0 {
		return defaults[0], nil
	} else {
		return strconv.Atoi(strv)
	}
}

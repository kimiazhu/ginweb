// Description: utils 签名模块,目前提供对用户中心签名模块的模块的支持。
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:08
package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/kimiazhu/log4go"
	"sort"
	"strings"
)

const (
	Secret_Key_User_Center = "xgsdkNiuxgsdkNiu"
	Secret_Key_Tako_Sdk    = "8NwwN2gitGTx4Nc9"
)

type SignError struct {
	Code int
	Msg  string
}

func (this SignError) Error() string {
	return this.Msg
}

func Sign(input, key string) string {
	keyForSign := []byte(key)
	h := hmac.New(sha1.New, keyForSign)
	h.Write([]byte(input))
	//	return base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 将Map格式化成待签名的字符串，
// 注意此方法会滤掉map中value为空的所有键值对
func formatForSign(data map[string]interface{}) string {
	if data == nil {
		log4go.Error("args could not be nil!!!")
		return ""
	}
	var validKeys []string
	for k := range data {
        switch value := data[k].(type) {
        case string:
            if value != "" {
                validKeys = append(validKeys, k)
            }
        case []interface{}:
            log4go.Finest("sign value is a slice, value=", value)
            if value != nil && len(value) > 0 {
                validKeys = append(validKeys, k)
            }
        }
	}

	sort.Strings(validKeys)

	var kv []string
	for _, k := range validKeys {
        switch value := data[k].(type) {
        case string:
            kv = append(kv, fmt.Sprintf("%s=%s", k, value))
        case []interface{}:
            vs := fmt.Sprintf("%s", value)
            kv = append(kv, fmt.Sprintf("%s=%s", k, strings.Replace(vs, " ", ",", -1)))
        }

	}
	dataForSign := strings.Join(kv, "&")
	log4go.Debug(fmt.Sprintf("data for sign: %s", dataForSign))
	return dataForSign
}

// 获取签名串，不包含其他任何东西，仅仅是签名字符串本身
func GetSignStr(data map[string]interface{}, key string) (string, error) {
	if data == nil || key == "" {
		log4go.Error("args could not be nil!!!")
		return "", SignError{Code: -1, Msg: "Arguments Error"}
	}
	dataForSign := formatForSign(data)
	log4go.Debug(fmt.Sprintf("data for sign: %s", dataForSign))
	return Sign(dataForSign, key), nil
}

// 构造已签名的请求字符串，
// 形式如： key1=val1&key2=val2&sign=signvalue
// 主要用于GET请求
func BuildSignedStr(data map[string]interface{}, key string) (string, error) {
	if data == nil || key == "" {
		log4go.Error("args could not be nil!!!")
		return "", SignError{Code: -1, Msg: "Arguments Error"}
	}
	dataForSign := formatForSign(data)
	log4go.Debug(fmt.Sprintf("data for sign: %s", dataForSign))
	sn := Sign(dataForSign, key)
	return dataForSign + "&sign=" + sn, nil
}

// 构造已签名的Json格式的字符串，主要用于POST请求的Body
func BuildSignedJsonStr(data map[string]interface{}, key string) (string, error) {
	if data == nil || key == "" {
		log4go.Error("args could not be nil!!!")
		return "", SignError{Code: -1, Msg: "Arguments Error"}
	}
	dataForSign := formatForSign(data)
	log4go.Debug(fmt.Sprintf("data for sign: %s", dataForSign))
	data["sign"] = Sign(dataForSign, key)
	js, jsErr := json.Marshal(data)
	if jsErr != nil {
		log4go.Error(fmt.Sprintf("format map to Json error"))
		return "", SignError{Code: -1, Msg: "Error format request to Json"}
	}
	return string(js), nil
}

// 签名校验，验证通过返回true，否则返回false
func SignVerify(data map[string]interface{}, secret string, sign string) bool {
	if data == nil || sign == "" {
		log4go.Error("args could not be nil!!!")
		return false
	}
	if data["sign"] != "" {
		// map中存在sign字段，需要将它移除
		delete(data, "sign")
	}
	signStr, _ := GetSignStr(data, secret)
    log4go.Info(fmt.Sprintf("verify result: orginSign=%s, correctSign=%s", sign, signStr))
	return signStr == sign
}

// 签名校验，原始签名串需要作为Json值放入到jsonStr参数中。
// 验证通过返回true，否则返回false
func SignVerify2(jsonStr string, secret string) bool {
	maps := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &maps)
	if err != nil {
		log4go.Error("args is not a json!", jsonStr)
		return false
	}

	if originSign := maps["sign"]; originSign == nil {
		log4go.Error("sign segment not found or is empty")
		return false
	} else {
	    return SignVerify(maps, secret, originSign.(string))
    }
}

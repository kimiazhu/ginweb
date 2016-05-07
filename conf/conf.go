// Description: conf 项目配置项，也用于读取当前目录下的app.conf中的项目
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:19
package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"reflect"
)

type Config struct {
	SERVER struct {
		RUNMODE string `yaml:runmode`
		PORT    string `yaml:"port"`
	}

	DATABASE struct {
		HOST     string `yaml:"host"`
		NAME     string `yaml:"name"`
		USER     string `yaml:"user"`
		PASSWORD string `yaml:"password"`
	}

	EXT map[string]interface{} `yaml:"ext,flow"`
}

// Ext will return the value of the EXT config, the keys is a string
// separated by DOT(.). If you provide a default value, this method
// will return the it while the key cannot be found. otherwise it
// will raise a panic!
func (c *Config) Ext(keys string, defaultVal... interface{}) (interface{}) {
	r, e := c.ExtSep(keys, ".")
	if e != nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(e)
		}
	} else {
		return r
	}
}

func (c *Config) ExtString(keys string, defaultVal... interface{}) (string) {
	return fmt.Sprintf("%v", c.Ext(keys, defaultVal...))
}

func (c *Config) ExtInt(keys string, defaultVal... interface{}) (int) {
	return c.Ext(keys, defaultVal...).(int)
}

func (c *Config) ExtInt8(keys string, defaultVal... interface{}) (int8) {
	return c.Ext(keys, defaultVal...).(int8)
}

func (c *Config) ExtInt16(keys string, defaultVal... interface{}) (int16) {
	return c.Ext(keys, defaultVal...).(int16)
}

func (c *Config) ExtInt32(keys string, defaultVal... interface{}) (int32) {
	return c.Ext(keys, defaultVal...).(int32)
}

func (c *Config) ExtInt64(keys string, defaultVal... interface{}) (int64) {
	return c.Ext(keys, defaultVal...).(int64)
}

func (c *Config) ExtBool(keys string, defaultVal... interface{}) (bool) {
	return c.Ext(keys, defaultVal...).(bool)
}

func (c *Config) ExtFloat64(keys string, defaultVal... interface{}) (float64) {
	return c.Ext(keys, defaultVal...).(float64)
}

func (c *Config) ExtFloat32(keys string, defaultVal... interface{}) (float32) {
	return c.Ext(keys, defaultVal...).(float32)
}

// Ext will return the value of the EXT config, the keys is separated
// by the given sep string.
func (c *Config) ExtSep(keys, sep string) (interface{}, error) {
	ks := strings.Split(keys, sep);
	var result interface{}
	var isFinal, success bool
	result = c.EXT
	for _, k := range ks {
		result, isFinal, success = find(result, k)
		if !success {
			return "", fmt.Errorf("no such key: %v", k)
		} else if isFinal {
			return result, nil
		}
	}
	return "", fmt.Errorf("not found")
}

func find(v interface{}, key interface{}) (result interface{}, isFinal, success bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		result = m[key.(string)]
		success = true
		isFinal = (reflect.TypeOf(result) != nil && reflect.TypeOf(result).Kind() != reflect.Map)
	case map[interface{}]interface{}:
		result = m[key]
		success = true
		isFinal = (reflect.TypeOf(result) != nil && reflect.TypeOf(result).Kind() != reflect.Map)
	}
	return
}

var (
	Conf = Config{}
)

func LoadConf(path string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	yaml.Unmarshal(c, &Conf)
}

func Ext(keys string, defVal... interface{}) interface {}{
	return Conf.Ext(keys, defVal...)
}

// ExtString is a shorcut for (c *Config) ExtString()
func ExtString(keys string, defVal... interface{}) string {
	return Conf.ExtString(keys, defVal...)
}

// ExtBool is a shorcut for (c *Config) ExtBool()
func ExtBool(keys string, defVal... interface{}) bool {
	return Conf.ExtBool(keys, defVal...)
}

// ExtInt is a shorcut for (c *Config) ExtInt()
func ExtInt(keys string, defVal... interface{}) int {
	return Conf.ExtInt(keys, defVal...)
}

// ExtInt8 is a shorcut for (c *Config) ExtInt8()
func ExtInt8(keys string, defVal... interface{}) int8 {
	return Conf.ExtInt8(keys, defVal...)
}

// ExtInt16 is a shorcut for (c *Config) ExtInt16()
func ExtInt16(keys string, defVal... interface{}) int16 {
	return Conf.ExtInt16(keys, defVal...)
}

// ExtInt32 is a shorcut for (c *Config) ExtInt32()
func ExtInt32(keys string, defVal... interface{}) int32 {
	return Conf.ExtInt32(keys, defVal...)
}

// ExtInt64 is a shorcut for (c *Config) ExtInt64()
func ExtInt64(keys string, defVal... interface{}) int64 {
	return Conf.ExtInt64(keys, defVal...)
}

// ExtFloat32 is a shorcut for (c *Config) ExtFloat32()
func ExtFloat32(keys string, defVal... interface{}) float32 {
	return Conf.ExtFloat32(keys, defVal...)
}

// ExtFloat64 is a shorcut for (c *Config) ExtFloat64()
func ExtFloat64(keys string, defVal... interface{}) float64 {
	return Conf.ExtFloat64(keys, defVal...)
}


func init() {
	LoadConf("conf.yml")
}

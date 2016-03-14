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
func (c *Config) Ext(keys string, defaultVal... string) (string) {
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

// Ext will return the value of the EXT config, the keys is separated
// by the given sep string.
func (c *Config) ExtSep(keys, sep string) (string, error) {
	ks := strings.Split(keys, sep);
	var result interface{}
	var isFinal, success bool
	result = c.EXT
	for _, k := range ks {
		result, isFinal, success = find(result, k)
		if !success {
			return "", fmt.Errorf("no such key: %v", k)
		} else if isFinal {
			return result.(string), nil
		}
	}
	return "", fmt.Errorf("not found")
}

func find(v interface{}, key interface{}) (result interface{}, isFinal, success bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		result = m[key.(string)]
		success = true
		_, isFinal = result.(string)
	case map[interface{}]interface{}:
		result = m[key]
		success = true
		_, isFinal = result.(string)
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

func ExtString(key string, defVal... string) (val string, err error) {
	if v := Conf.EXT[key]; v != "" {
		val = fmt.Sprintf("%v", v)
		return
	}

	if len(defVal) > 0 {
		val = defVal[0]
		return
	}

	err = fmt.Errorf("no such key[%s] and no default value provided", key)
	return
}

func init() {
	LoadConf("conf.yml")
}

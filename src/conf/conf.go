// Description: conf 项目配置项，也用于读取当前目录下的app.conf中的项目
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:19
package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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

var (
	Conf = Config{}
)

func init() {
	c, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	yaml.Unmarshal(c, &Conf)
}

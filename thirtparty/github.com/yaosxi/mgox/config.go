package mgox

import (
	"bufio"
	"github.com/kimiazhu/log4go"
	"io"
	"os"
	"strings"
	"fmt"
)

type dbconfig struct {
	Host     string
	Database string
	Username string
	Password string
}

var DBConfig dbconfig

type PropertyReader struct {
	m map[string]string
}

func (p *PropertyReader) init(path string) {

	p.m = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))

		//log.Println(s)

		if strings.Index(s, "#") == 0 {
			continue
		}

		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		p.m[frist] = strings.TrimSpace(second)
	}
}

func (c PropertyReader) Read(key string) string {
	v, found := c.m[key]
	if !found {
		return ""
	}
	return v
}

func LoadConfig(path string) {
	db := new(PropertyReader)
	db.init(path)
	DBConfig.Host = db.m["host"]
	DBConfig.Database = db.m["database"]
	DBConfig.Username = db.m["username"]
	DBConfig.Password = db.m["password"]
	log4go.Debug(fmt.Sprintf("host=%s,database=%s,username=%s\n", DBConfig.Host, DBConfig.Database, DBConfig.Username))
}

func init() {
	LoadConfig("conf/mgox.properties")
}

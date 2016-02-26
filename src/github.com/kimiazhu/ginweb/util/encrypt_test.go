// Description: utils
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:08
package utils

import (
	"testing"
	log "github.com/kimiazhu/log4go"
)

func init() {
//	log.LoadConfiguration("../../../conf/log4go.xml")
}

func TestCBCEncrypterAESNoPadding(t *testing.T) {
	log.LoadConfiguration("../../../conf/log4go.xml")
	r, _ := CBCEncrypterAESNoPadding("sdkportal__490ae2ab47e032841a86a5c7bed04e81db1147b0b6351b3d62268cbbd149bea2", Secret_Key_User_Center, Secret_Key_User_Center)
	t.Log("\nresult is:", r)
}
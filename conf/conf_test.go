// Copyright 2011 ZHU HAIHUA <kimiazhu@gmail.com>. 
// All rights reserved.
// Use of this source code is governed by MIT 
// license that can be found in the LICENSE file.

// Description: conf
// Author: ZHU HAIHUA
// Since: 9/19/16
package conf

import (
	"testing"
)

func TestDatabase_Ext(t *testing.T) {
	LoadConf("testdata/conf.yml")
	maxIdle := Conf.DATABASE.Ext("maxIdle").(int)
	if maxIdle != 10 {
		t.Errorf("maxIdle need to be 10, but now it is: %v", maxIdle)
	}
}

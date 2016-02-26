// Created by yaoshuangxi

package mgox

import (
	"github.com/kimiazhu/log4go"
	"gopkg.in/mgo.v2"
	"sync"
)

var (
	dbSession *mgo.Session
	lock = sync.Mutex{}
)


func GetDatabase() (*mgo.Database, error) {
	lock.Lock()
	defer lock.Unlock()
	if dbSession == nil {
		var err error
		log4go.Info("Dial to %s", DBConfig.Host)
		dbSession, err = mgo.Dial(DBConfig.Host)
		if err != nil {
			return nil, err
		}
		dbSession.SetMode(mgo.Monotonic, true)
	}

	log4go.Finest("Try to open DB connnection: %s", DBConfig.Database)
	database := dbSession.Clone().DB(DBConfig.Database)

	if DBConfig.Username != "" {
		loginErr := database.Login(DBConfig.Username, DBConfig.Password)
		if loginErr != nil {
			return nil, loginErr
		}
	}

	log4go.Finest("Opened DB connnection successfully")

	return database, nil
}

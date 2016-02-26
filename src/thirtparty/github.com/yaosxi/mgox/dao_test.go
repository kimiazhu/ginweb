package mgox_test

import (
	"github.com/yaosxi/mgox"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
	"fmt"
)

func handleError(t *testing.T, err error) bool {
	if err == nil {
		return false
	}
	panic(err)
	t.Fail()
	return true
}

type Log4go struct {

}

func (log Log4go)Debug(v interface{})  {
	fmt.Println(v)
}

func (log Log4go)Critical(v interface{})  {
	fmt.Println("*************** Error ***************")
	fmt.Println(v)
	fmt.Println("*************** Error ***************")
}

var log4go Log4go

var UserCollectionName = new(User)

type User struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name         string        `json:"name"`
	Age          int           `json:"age"`
	Sex          int           `json:"sex"`
	FirstCreator string        `json:"firstcreator"`
	FirstCreated time.Time     `json:"firstcreated"`
	LastModifier string        `json:"lastmodifier"`
	LastModified time.Time     `json:"lastmodified"`
}

type IpCache struct {
	Ip           string `json:"ip" bson:"_id,omitempty"`
	Country         string        `json:"country"`
}

func getFirst() User {
	var user User
	mgox.New().Find().First(&user)
	return user
}

func TestConnection(t *testing.T) {
	dao := mgox.Connect()
	defer dao.Close()
	var user = new(User)
	dao.Find().IgnoreNFE().First(user)
	log4go.Debug(user.Id)
	dao.Find().IgnoreNFE().First(user)
	err := dao.Find().IgnoreNFE().First(user)
	if handleError(t, err) {
		return
	}
}

func TestInsert(t *testing.T) {
	err := mgox.New("111111").Insert(
		&User{Name: "yaosxi", Age: 2, Sex: 1},
		&User{Name: "yaosxi2", Age: 3, Sex: 1},
		bson.M{"name": "yaosxi3", "age": 3, "sex": 1},
	)
	if handleError(t, err) {
		return
	}

	err = mgox.New("111111").Insert(
		User{Name: "yaosxi4", Age: 2, Sex: 1},
	)
	if handleError(t, err) {
		return
	}
}

func TestRemove(t *testing.T) {

	user := getFirst()
	err := mgox.New().Remove(user, user.Id)
	if handleError(t, err) {
		return
	}

	user = getFirst()
	err = mgox.New().Remove("user", user.Id)
	if handleError(t, err) {
		return
	}

	err = mgox.New().Remove(user, bson.M{"name": "carson2"})
	if handleError(t, err) {
		return
	}
}

func TestSet(t *testing.T) {

	user := getFirst()

	name := "ysx111"
	err := mgox.New().Set(user, user.Id, "name", name)
	if handleError(t, err) {
		return
	}
	user = getFirst()
	if user.Name != name {
		t.Fail()
		return
	}

	name = "ysx222"
	err = mgox.New().Set(&user, user.Id, "name", name)
	if handleError(t, err) {
		return
	}
	user = getFirst()
	if user.Name != name {
		t.Fail()
		return
	}

	name = "ysx333"
	err = mgox.New().Set("user", user.Id, "name", name)
	if handleError(t, err) {
		return
	}
	user = getFirst()
	if user.Name != name {
		t.Fail()
	}

	name = "ysx444"
	err = mgox.New().Set(UserCollectionName, user.Id, bson.M{"name": name})
	if handleError(t, err) {
		return
	}
	user = getFirst()
	if user.Name != name {
		t.Fail()
	}
}

func TestInc(t *testing.T) {
	user := getFirst()
	err := mgox.New("uuuuuuuuuuuuuu").Inc(user, user.Id, "age", 1)
	if handleError(t, err) {
		return
	}
	age := user.Age
	user = getFirst()
	if user.Age != age+1 {
		t.Fail()
	}
}

func TestReplace(t *testing.T) {
	user := getFirst()
	err := mgox.New().Replace(user, user.Id, "name", "carson", "age", 10)
	if handleError(t, err) {
		return
	}
	user = getFirst()
	if user.Name != "carson" {
		log4go.Critical(user.Name)
		t.Fail()
		return
	}
	err = mgox.New().ReplaceDoc(User{Id: user.Id, Name: "carson2"})
	if handleError(t, err) {
		return
	}
	user = getFirst()
	if user.Name != "carson2" {
		log4go.Critical(user.Name)
		t.Fail()
	}
}

func TestReplaceDoc(t *testing.T) {

	ip := "127.0.0.1"

	//err := mgox.New("111111").Insert(
	//	&IpCache{Ip:ip, Country: "CN"},
	//)
	//if handleError(t, err) {
	//	return
	//}

	err := mgox.New().ReplaceDocById(ip, IpCache{Ip: ip, Country: "CN"})
	if handleError(t, err) {
		return
	}

}

func TestCount(t *testing.T) {
	n, err := mgox.New().Find().Count(UserCollectionName)
	if handleError(t, err) {
		return
	}
	if n != 4 {
		log4go.Debug(n)
		t.Fail()
	}
}

func TestFirst(t *testing.T) {
	var user User
	err := mgox.New().Find().IgnoreNFE().First(&user)
	if handleError(t, err) {
		return
	}
	err = mgox.New().Find().First(&user)
	if handleError(t, err) {
		return
	}
}

func TestLast(t *testing.T) {
	var user User
	err := mgox.New().Find("name", "yaosxi").IgnoreNFE().Last(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user.Name)

	user = User{}
	err = mgox.New().Find("name", "name").IgnoreNFE().Last(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user.Name)
}

func TestFind(t *testing.T) {
	var users []User
	err := mgox.New().Find().Result(&users)
	if handleError(t, err) {
		return
	}
	log4go.Debug(users)

	p := mgox.Page{Count: 1}
	err = mgox.New().Find().Page(&p).Sort("age", "-name").Result(&users)
	if handleError(t, err) {
		return
	}
	log4go.Debug(users)

	var user User
	err = mgox.New().Find(users[0].Id).Result(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user)

	user = User{}
	err = mgox.New().Find("name", "yaosxi2").Result(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user)
}

func TestGet(t *testing.T) {

	var user User
	err := mgox.New().Get().Result(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user)

	id := user.Id
	user = User{}
	err = mgox.New().Get(id).Result(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user)

	user = User{}
	err = mgox.New().Get("name", "yaosxi2").Result(&user)
	if handleError(t, err) {
		return
	}
	log4go.Debug(user)

	var users []User
	err = mgox.New().Find().Result(&users)
	if handleError(t, err) {
		return
	}
}

func TestSelect(t *testing.T) {

	var users []User
	err := mgox.New().Find().Select("name").Result(&users)
	if handleError(t, err) {
		return
	}
	log4go.Debug("-------------------------------")
	log4go.Debug(users)
	log4go.Debug("-------------------------------")
}

func TestExist(t *testing.T) {

	b, err := mgox.New().Find().Exist("user")
	if handleError(t, err) {
		return
	}
	if !b {
		log4go.Debug(b)
		t.Fail()
		return
	}

	b, err = mgox.New().Find(bson.ObjectIdHex("56597ab9f918ad09b4000001")).Exist("user")
	if handleError(t, err) {
		return
	}
	if !b {
		log4go.Debug(b)
		t.Fail()
		return
	}

	b, err = mgox.New().Find("name", "yaosxi").Exist("user")
	if handleError(t, err) {
		return
	}
	if !b {
		log4go.Debug(b)
		t.Fail()
		return
	}

	b, err = mgox.New().Find(bson.M{"name": "yaosxi"}).Exist("user")
	if handleError(t, err) {
		return
	}
	if !b {
		log4go.Debug(b)
		t.Fail()
		return
	}
}
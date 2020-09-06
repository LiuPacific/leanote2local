package tmig

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"regexp"
	"testing"
	"time"
)

func err_handler(err error) {
	fmt.Printf("err_handler, error:%s\n", err.Error())
	panic(err.Error())
}
func GetSession() *mgo.Session {
	dail_info := &mgo.DialInfo{
		Addrs:    []string{"127.0.0.1"},
		Direct:   false,
		Timeout:  time.Second * 1,
		Database: "ir",
		//Source:    "admin",
		//Username:  "test1",
		//Password:  "123456",
		PoolLimit: 1024,
	}
	session, err := mgo.DialWithInfo(dail_info)
	if err != nil {
		fmt.Printf("mgo dail error[%s]\n", err.Error())
		err_handler(err)
	}

	// set mode
	session.SetMode(mgo.Monotonic, true)

	return session
}

func TestCreateTopDir(t *testing.T) {
	session := GetSession()
	defer session.Close()
	err := CreateTopDir(session)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestT(t *testing.T) {
	str0 := "asdf![title](/api/file/getImage?fileId=5f1d3aba4b6c3f6fa6000e49)dasfdfasdf![title](/api/file/getImage?fileId=5wfaaba4b6c3f6fa6000e49)dasfdf"
	regStr0 := "!\\[.*?\\]\\(/api/file/getImage\\?fileId=(\\w*?)\\)"
	regCompiled0 := regexp.MustCompile(regStr0)
	submatch0 := regCompiled0.FindAllStringSubmatch(str0, -1)
	fmt.Println(submatch0)
	fmt.Println("---")
	for i:=0;i<len(submatch0);i++  {
		fmt.Println(submatch0[i][1])
	}
}

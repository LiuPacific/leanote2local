package main

import (
	"fmt"
	"leanote2vnote/note/tmig"
	"time"

	"gopkg.in/mgo.v2"
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

func main() {
	session := GetSession()
	defer session.Close()
	err := tmig.CreateTopDir(session)
	if err != nil {
		fmt.Println(err.Error())
	}
}

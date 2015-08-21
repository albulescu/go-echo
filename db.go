package main

import (
	"gopkg.in/mgo.v2"
	"strconv"
	"time"
	//"gopkg.in/mgo.v2/bson"
)

type MgoProxy struct {
	sess *mgo.Session
	db   *mgo.Database
}

func (m *MgoProxy) Close() {
	m.sess.Close()
}

func (m *MgoProxy) Init() {

	timeout, err := time.ParseDuration(configDb.Timeout)

	if err != nil {
		die("Invalid database timeout")
	}

	info := &mgo.DialInfo{
		Addrs:    []string{configDb.Host},
		Timeout:  timeout,
		Database: configDb.Database,
		Username: configDb.Username,
		Password: configDb.Password,
		Source:   configDb.Source,
	}

	sess, err := mgo.DialWithInfo(info)

	if err != nil {
		die("Fail to connect to mongo database %s:%s@%s:%s", configDb.Username,
			configDb.Password,
			configDb.Host,
			strconv.Itoa(configDb.Port))
	}

	sess.SetMode(mgo.Monotonic, true)

	m.sess = sess
	m.db = sess.DB(configDb.Database)
}

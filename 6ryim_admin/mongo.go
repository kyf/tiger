package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	host    string
	port    string
	user    string
	pass    string
	dbname  string
	session *mgo.Session
	db      *mgo.Database
}

var (
	mongodbServer   string = "127.0.0.1"
	mongodbPort     string = "27017"
	mongodbName     string = "call_center"
	mongodbUser     string = ""
	mongodbPass     string = ""
	mongodbPoolSize int    = 300
)

func NewMongoClient() *Mongo {
	return &Mongo{
		host:   mongodbServer,
		port:   mongodbPort,
		user:   mongodbUser,
		pass:   mongodbPass,
		dbname: mongodbName,
	}
}

func (this *Mongo) Connect() error {
	url := fmt.Sprintf("%s:%s", this.host, this.port)
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetPoolLimit(mongodbPoolSize)
	this.session = session
	this.db = session.DB(this.dbname)
	//err = this.db.Login(this.user, this.pass)
	return err
}

func (this *Mongo) Add(coll string, data ...interface{}) error {
	return this.db.C(coll).Insert(data...)
}

func (this *Mongo) Find(coll string, query interface{}) *mgo.Query {
	return this.db.C(coll).Find(query)
}

func (this *Mongo) Update(coll string, query interface{}, data interface{}) error {
	return this.db.C(coll).Update(query, data)
}

func (this *Mongo) Remove(coll string, query interface{}) error {
	return this.db.C(coll).Remove(query)
}

func (this *Mongo) Close() {
	this.session.Close()
}

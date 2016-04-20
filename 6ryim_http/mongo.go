package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	//"github.com/go-mgo/mgo/bson"
)

type Mongo struct {
	host      string
	port      string
	user      string
	pass      string
	dbname    string
	msgdbname string
	session   *mgo.Session
	db        *mgo.Database
	msgdb     *mgo.Database
}

var (
	mongodbServer   string
	mongodbPort     string
	mongodbName     string
	msgMongodbName  string
	mongodbUser     string
	mongodbPass     string
	mongodbPoolSize int
)

func InitMongodb(pmongodbServer, pmongodbPort, pmongodbUser, pmongodbName, pMsgMongodbName, pmongodbPass string, pmongodbPoolSize int) {
	mongodbServer = pmongodbServer
	mongodbPort = pmongodbPort
	mongodbName = pmongodbName
	msgMongodbName = pMsgMongodbName
	mongodbUser = pmongodbUser
	mongodbPass = pmongodbPass
	mongodbPoolSize = pmongodbPoolSize

}

func NewMongoClient() *Mongo {
	return &Mongo{
		host:      mongodbServer,
		port:      mongodbPort,
		user:      mongodbUser,
		pass:      mongodbPass,
		dbname:    mongodbName,
		msgdbname: msgMongodbName,
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
	this.msgdb = session.DB(this.msgdbname)
	err = this.db.Login(this.user, this.pass)
	return err
	//return nil
}

func (this *Mongo) Add(coll string, data ...interface{}) error {
	return this.msgdb.C(coll).Insert(data...)
}

func (this *Mongo) Find(coll string, query interface{}) *mgo.Query {
	return this.msgdb.C(coll).Find(query)
}

func (this *Mongo) Remove(coll string, query interface{}) error {
	return this.msgdb.C(coll).Remove(query)
}

func (this *Mongo) Close() {
	this.session.Close()
}

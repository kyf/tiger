package im_type

import (
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Message    string `json:"message"`
	OrderId    string `json:"orderid"`
	FromType   string `json:"fromtype"`
	ToType     string `json:"totype"`
	MsgType    string `json:"msgtype"`
	CreateTime string `json:"createtime"`
	IsSystem   string `json:"issystem"`
	SystemType string `json:"systemtype"`
	Source     string `json:"source"`
}

type BsonMessage struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	From       string        `bson:"from" json:"from"`
	To         string        `bson:"to" json:"to"`
	Message    string        `bson:"message" json:"message"`
	OrderId    string        `bson:"orderid" json:"orderid"`
	FromType   string        `bson:"fromtype" json:"fromtype"`
	ToType     string        `bson:"totype" json:"totype"`
	MsgType    string        `bson:"msgtype" json:"msgtype"`
	CreateTime string        `bson:"createtime" json:"createtime"`
	IsSystem   string        `bson:"issystem" json:"issystem"`
	SystemType string        `bson:"systemtype" json:"systemtype"`
	Source     string        `bson:"source" json:"source"`
}

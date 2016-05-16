package main

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	FAST_REPLY_TABLE string = "cc_fastreply"

	AUTO_REPLY_TABLE string = "cc_autoreply"
)

type FastReply struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Content    string        `json:"content" bson:"content"`
	Author     string        `json:"author" bson:"author"`
	CreateTime int64         `json:"createtime" bson:"createtime"`
}

func (fr *FastReply) Add(mgo *Mongo) error {
	return mgo.Add(FAST_REPLY_TABLE, fr)
}

func (fr *FastReply) Update(mgo *Mongo) error {
	query := bson.M{
		"_id": fr.Id,
	}
	data := bson.M{
		"content": fr.Content,
	}
	return mgo.Update(FAST_REPLY_TABLE, query, data)
}

func FastReplyList(mgo *Mongo) ([]FastReply, error) {
	var result []FastReply
	err := mgo.Find(FAST_REPLY_TABLE, nil).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (fr *FastReply) Remove(mgo *Mongo) error {
	data := bson.M{
		"_id": fr.Id,
	}
	return mgo.Remove(FAST_REPLY_TABLE, data)
}

type AutoReply struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	Content    string        `json:"content" bson:"content"`
	FromHour   int           `json:"fromhour" bson:"fromhour"`
	FromMinute int           `json:"fromminute" bson:"fromminute"`
	ToHour     int           `json:"tohour" bson:"tohour"`
	ToMinute   int           `json:"tominute" bson:"tominute"`
}

func (ar *AutoReply) Add(mgo *Mongo) error {
	return mgo.Add(AUTO_REPLY_TABLE, ar)
}

func (ar *AutoReply) Update(mgo *Mongo) error {
	query := bson.M{
		"_id": ar.Id,
	}
	data := bson.M{
		"content":    ar.Content,
		"fromhour":   ar.FromHour,
		"fromminute": ar.FromMinute,
		"tohour":     ar.ToHour,
		"tominute":   ar.ToMinute,
	}
	return mgo.Update(AUTO_REPLY_TABLE, query, data)
}

func AutoReplyList(mgo *Mongo) ([]AutoReply, error) {
	var result []AutoReply
	err := mgo.Find(AUTO_REPLY_TABLE, nil).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ar *AutoReply) Remove(mgo *Mongo) error {
	data := bson.M{
		"_id": ar.Id,
	}
	return mgo.Remove(AUTO_REPLY_TABLE, data)
}

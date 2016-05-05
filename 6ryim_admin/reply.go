package main

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	FAST_REPLY_TABLE string = "cc_fastreply"
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

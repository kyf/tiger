package main

import (
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MessageType int

type Message struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Openid   string        `json:"openid" bson:"openid"`
	Created  int64         `json:"ts" bson:"created"`
	Content  string        `json:"content" bson:"content"`
	MsgType  MessageType   `json:"msgType" bson:"msgtype"`
	Opid     string        `json:"opid" bson:"opid"`
	Fromtype int           `json:"fromtype" bson:"fromtype"`
	Source   int           `json:"source" bson:"source"`
}

const (
	CLIENT_DEADLINE time.Duration = time.Second * 60 * 5 //5 minutes

	MSG_TYPE_TEXT  MessageType = 1
	MSG_TYPE_IMAGE MessageType = 2
	MSG_TYPE_AUDIO MessageType = 3

	MSG_FROM_TYPE_USER int = 1
	MSG_FROM_TYPE_OP   int = 2

	MSG_SOURCE_WX      int = 1
	MSG_SOURCE_IOS     int = 2
	MSG_SOURCE_ANDROID int = 3
	MSG_SOURCE_PC      int = 4

	CC_MESSAGE_TABLE = "cc_message"
)

func listMessage(openid string, mgo *Mongo) ([]Message, error) {
	var result []Message
	err := mgo.Find(CC_MESSAGE_TABLE, bson.M{"openid": openid}).Sort("-_id").Limit(100).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func storeMessage(msg Message, mgo *Mongo) error {
	if msg.Source == 0 {
		msg.Source = MSG_SOURCE_WX
	}
	data := bson.M{
		"_id":      bson.NewObjectId(),
		"openid":   msg.Openid,
		"created":  msg.Created,
		"content":  msg.Content,
		"msgtype":  msg.MsgType,
		"opid":     msg.Opid,
		"fromtype": msg.Fromtype,
		"source":   msg.Source,
	}
	err := mgo.Add(CC_MESSAGE_TABLE, data)
	if err != nil {
		return err
	}
	return nil

}

type WaitList struct {
	waitPool map[string][]Message
	locker   sync.Mutex
}

func NewWaitList() *WaitList {
	return &WaitList{waitPool: make(map[string][]Message)}
}

func (wl *WaitList) Add(msg Message) {
	wl.locker.Lock()
	if _, ok := wl.waitPool[msg.Openid]; !ok {
		wl.waitPool[msg.Openid] = make([]Message, 0)
	}
	wl.waitPool[msg.Openid] = append(wl.waitPool[msg.Openid], msg)
	wl.locker.Unlock()
}

func (wl *WaitList) Fetch(opid, openid string) bool {
	wl.locker.Lock()
	defer wl.locker.Unlock()
	var msgs []Message
	var ok bool
	if msgs, ok = wl.waitPool[openid]; !ok {
		msgs = make([]Message, 0)
	}
	if status := defaultOL.bind(opid, openid, msgs); !status {
		return false
	}
	delete(wl.waitPool, openid)
	return true
}

var defaultWL *WaitList = NewWaitList()

var defaultOL *Online = NewOnline()

func rmExpire() {
	defaultOL.poolLocker.Lock()
	for opid, clients := range defaultOL.olPool {
		tmpClients := make([]Client, 0, len(clients))
		for _, client := range clients {
			if client.lastTS.Add(CLIENT_DEADLINE).Unix() > time.Now().Unix() {
				tmpClients = append(tmpClients, client)
			} else {
				delete(defaultOL.userMapping, client.openid)
			}
		}

		if len(tmpClients) == 0 {
			delete(defaultOL.olPool, opid)
		} else {
			defaultOL.olPool[opid] = tmpClients
		}
	}
	defaultOL.poolLocker.Unlock()

}

func init() {
	ticker := time.NewTicker(CLIENT_DEADLINE)
	go func() {
		for {
			select {
			case <-ticker.C:
				rmExpire()
			}
		}
	}()
}

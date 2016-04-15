package main

import (
	"sync"
	"time"
)

type MessageType int

type Message struct {
	openid  string `json:"openid"`
	created int64
	content string `json:"content"`
	msgType MessageType
}

const (
	CLIENT_DEADLINE time.Duration = time.Second * 60 * 5 //5 minutes

	MSG_TYPE_TEXT  MessageType = 1
	MSG_TYPE_IMAGE MessageType = 2
	MSG_TYPE_AUDIO MessageType = 3
)

type WaitList struct {
	waitPool map[string]Message
	locker   sync.Mutex
}

func NewWaitList() *WaitList {
	return &WaitList{waitPool: make(map[string]Message)}
}

func (wl *WaitList) Add(msg Message) {
	wl.locker.Lock()
	wl.waitPool[msg.openid] = msg
	wl.locker.Unlock()
}

func (wl *WaitList) Fetch(opid, openid string) bool {
	wl.locker.Lock()
	var msg Message
	var ok bool
	if msg, ok = wl.waitPool[openid]; !ok {
		return false
	}
	defaultOL.bind(opid, openid, msg)
	delete(wl.waitPool, openid)
	wl.locker.Unlock()
	return true
}

var defaultWL *WaitList = NewWaitList()

var defaultOL *Online = NewOnline()

func rmExpire() {
	defaultOL.poolLocker.Lock()
	for opid, clients := range defaultOL.olPool {
		tmpClients := make([]Client, 0, len(clients))
		for _, client := range clients {
			if client.lastTS.Add(CLIENT_DEADLINE).Unix() < time.Now().Unix() {
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

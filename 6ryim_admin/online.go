package main

import (
	"strings"
	"sync"
	"time"
)

type Client struct {
	openid    string `json:"openid"`
	lastMsg   *Message
	unRead    []Message
	msgLocker sync.Mutex
	lastTS    time.Time
	tsLocker  sync.Mutex
}

type Online struct {
	olPool     map[string][]Client
	poolLocker sync.Mutex

	userMapping map[string]string
}

func NewOnline() *Online {
	return &Online{olPool: make(map[string][]Client), userMapping: make(map[string]string)}
}

func (ol *Online) getClientByOpenid(openid string) *Client {
	ol.poolLocker.Lock()
	defer ol.poolLocker.Unlock()
	if opid, ok := ol.userMapping[openid]; ok {
		if clients, ok := ol.olPool[opid]; ok {
			for index, client := range clients {
				if strings.EqualFold(client.openid, openid) {
					return &(ol.olPool[opid][index])
				}
			}
		}
	}

	return nil
}

func (ol *Online) findOpByUser(openid string) string {
	ol.poolLocker.Lock()
	defer ol.poolLocker.Unlock()
	if opid, ok := ol.userMapping[openid]; ok {
		return opid
	} else {
		return ""
	}
}

func (ol *Online) bind(opid, openid string, msgs []Message) bool {
	ol.poolLocker.Lock()
	defer ol.poolLocker.Unlock()
	if _, ok := ol.userMapping[openid]; ok {
		return false
	}
	ol.userMapping[openid] = opid
	lastMsg := new(Message)
	if len(msgs) > 0 {
		lastMsg = &(msgs[len(msgs)-1])
	}
	ol.olPool[opid] = append(ol.olPool[opid], Client{openid: openid, lastTS: time.Now(), lastMsg: lastMsg, unRead: msgs})
	return true
}

func (ol *Online) unbind(opid, openid string) {
	ol.poolLocker.Lock()
	defer ol.poolLocker.Unlock()
	if clients, ok := ol.olPool[opid]; ok {
		var currentIndex int = -1
		for index, client := range clients {
			if strings.EqualFold(openid, client.openid) {
				currentIndex = index
				break
			}
		}

		if currentIndex > -1 {
			if len(clients) == 1 {
				delete(ol.olPool, opid)
				return
			}

			if currentIndex == len(clients)-1 {
				ol.olPool[opid] = ol.olPool[opid][:len(clients)-1]
			} else {
				ol.olPool[opid][currentIndex] = ol.olPool[opid][len(clients)-1]
				ol.olPool[opid] = ol.olPool[opid][:len(clients)-1]
			}
		}

	}
}

func (ol *Online) getClient(opid, openid string) *Client {
	ol.poolLocker.Lock()
	defer ol.poolLocker.Unlock()
	if clients, ok := ol.olPool[opid]; ok {
		for _, client := range clients {
			if strings.EqualFold(openid, client.openid) {
				return &client
			}
		}
	}

	return nil
}

func (client *Client) appendMsg(msg Message) {
	client.msgLocker.Lock()
	client.lastMsg = &msg
	client.unRead = append(client.unRead, msg)
	client.msgLocker.Unlock()
}

func (client *Client) fetchMsg() map[string]interface{} {
	client.msgLocker.Lock()
	last := client.lastMsg
	unRead := client.unRead
	client.unRead = client.unRead[:0]
	client.msgLocker.Unlock()
	return map[string]interface{}{"last": *last, "unread": unRead}
}

func (client *Client) refresh() {
	client.tsLocker.Lock()
	client.lastTS = time.Now()
	client.tsLocker.Unlock()
}

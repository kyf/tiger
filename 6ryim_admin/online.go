package main

import (
	"sync"
)

type Client struct {
	openid    string `json:"openid"`
	msgPool   []Message
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

func (ol *Online) findOpByUser(openid string) string {
	if opid, ok := ol.userMapping[openid]; ok {
		return opid
	} else {
		return ""
	}
}

func (ol *Online) bind(opid, openid string) {
	ol.poolLocker.Lock()
	ol.userMapping[openid] = opid
	ol.olPool[opid] = append(ol.olPool[opid], &Client{openid: openid, lastTS: time.Now(), msgPool: make([]Message, 0, 10)})
	ol.poolLocker.Unlock()
}

func (ol *Online) unbind(opid, openid string) {
	ol.poolLocker.Lock()
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
	ol.poolLocker.Unlock()
}

func (client *Client) appendMsg(msg Message) {
	client.msgLocker.Lock()
	client.msgPool = append(client.msgPool, msg)
	client.msgLocker.Unlock()
}

func (client *Client) fetchMsg() []Message {
	client.msgLocker.Lock()
	result := client.msgPool
	client.msgPool = client.msgPool[:0]
	client.msgLocker.Unlock()
	return result
}

func (client *Client) refresh() {
	client.tsLocker.Lock()
	client.lastTS = time.Now()
	client.tsLocker.Unlock()
}

package main

import (
	"encoding/json"
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

func newMsg(m []byte) (*Message, error) {
	var result Message
	err := json.Unmarshal(m, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

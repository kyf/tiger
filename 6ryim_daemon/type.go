package main

import (
	"encoding/json"

	im_type "github.com/kyf/6ryim/6ryim_http/im_type"
)

type Message im_type.Message

func newMsg(m []byte) (*Message, error) {
	var result Message
	err := json.Unmarshal(m, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

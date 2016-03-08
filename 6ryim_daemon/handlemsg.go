package main

import (
	"net/http"
)

func handleMsg(msg Message) {
	switch msg.ToType {
	case TERMINAL_ADMIN:
	case TERMINAL_IOS:
	case TERMINAL_ANDROID:
	case TERMINAL_WX:
	case TERMINAL_360STREAM:
	default:
		return false
	}
	return true
}

package main

import (
	"log"
)

func handleMsg(msg Message, logger *log.Logger) bool {
	//getMedia from weixin resource
	if strings.EqualFold(msg.ToType, TERMINAL_ADMIN) {
		if strings.EqualFold(msg.Source, MSG_SOURCE_WX) {
			switch msg.MsgType {
			case MSG_TYPE_IMAGE:
				fallthrough
			case MSG_TYPE_AUDIO:

			}
		}
	}

	if h.isOnline(msg.To) {
		return false
	}

	var err error = nil
	err = storeOffline(msg)
	if err != nil {
		logger.Printf("storeOffline err:%v", err)
	}

	if strings.EqualFold(msg.ToType, TERMINAL_ADMIN) {
		err = msg.sendAdminTpl()
	} else {

		switch msg.Source {
		case MSG_SOURCE_WX:
			err = msg.sendUserWX()
		case MSG_SOURCE_IOS:
			err = msg.sendUserIOS()
		case MSG_SOURCE_ANDROID:
			err = msg.sendUserAndroid()
		case MSG_SOURCE_360STREAM:
			err = msg.sendUser360Stream()
		default:
			return false
		}
	}

	if err != nil {
		logger.Printf("handleMsg err:%v", err)
	}
	return true
}

func serveMsgReceive() {

}

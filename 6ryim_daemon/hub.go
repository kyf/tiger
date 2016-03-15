package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type hub struct {
	online map[string]*connection

	olmutex sync.RWMutex

	message chan []byte

	register chan *connection

	unregister chan *connection
}

var h = hub{
	message:    make(chan []byte),
	register:   make(chan *connection),
	unregister: make(chan *connection),
	online:     make(map[string]*connection),
}

func (h *hub) isOnline(to string) bool {
	h.olmutex.RLock()
	defer h.olmutex.RUnlock()
	var result bool = false
	if _, ok := h.online[to]; ok {
		result = true
	}
	return result
}

func listOnline(w http.ResponseWriter) {
	h.olmutex.RLock()
	defer h.olmutex.RUnlock()
	result := make([]string, 0)
	for _, it := range h.online {
		result = append(result, it.tk)
	}

	response(w, true, "success", result)
}

func (h *hub) run(logger *log.Logger) {
	for {
		select {
		case c := <-h.register:
			logger.Printf("register %v", c)
			h.olmutex.Lock()
			if ct, ok := h.online[c.token]; ok {
				go func() {
					lastc := ct
					lastc.send <- []byte("you are kicked out!")
					time.Sleep(time.Second * 5)
					h.unregister <- lastc
				}()
			}
			h.online[c.token] = c
			h.olmutex.Unlock()

			mli, err := fetchOffline(c.token)
			if err == nil {
				for _, mit := range mli {
					c.send <- []byte(mit)
				}
			}

		case c := <-h.unregister:
			logger.Printf("unregister %v", c)
			h.olmutex.Lock()
			if _, ok := h.online[c.token]; ok {
				delete(h.online, c.token)
				close(c.send)
			}
			h.olmutex.Unlock()
		case m := <-h.message:
			msg, err := newMsg(m)
			logger.Printf("receive message %v", msg)
			if err != nil {
				logger.Printf("newMsg err:%v", err)
				break
			}

			logger.Printf("send message 1 ...")
			to := msg.To
			msg.CreateTime = getFormatNow("num")
			var status bool
			logger.Printf("send message 2 ...")
			status, msg = handleMsg(msg, logger)
			logger.Printf("send message 3 ...")
			go func() {
				err := storeMessage(*msg)
				if err != nil {
					logger.Println("storeMessage err:%v", err)
				}

			}()

			logger.Printf("break before ...")
			if status {
				break
			}
			logger.Printf("break after ...")

			if c, ok := h.online[to]; ok {
				m, err = json.Marshal(msg)
				if err != nil {
					logger.Println("sendMessage json err:%v", err)
					break
				}
				select {
				case c.send <- m:
				default:
					h.olmutex.Lock()
					close(c.send)
					delete(h.online, to)
					var err error = nil
					err = storeOffline(*msg)
					if err != nil {
						logger.Printf("storeOffline err:%v", err)
					}
					h.olmutex.Unlock()
				}
			}
		}
	}

	logger.Printf("exit hub run ...")
}

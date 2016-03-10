package main

import (
	"log"
)

type hub struct {
	online map[string][]*connection

	message chan []byte

	register chan *connection

	unregister chan *connection
}

var h = hub{
	message:    make(chan []byte),
	register:   make(chan *connection),
	unregister: make(chan *connection),
	online:     make(map[string][]*connection),
}

func (h *hub) isOnline(to string) bool {
	var result bool = false
	if li, ok := h.online[to]; ok {
		if len(li) > 0 {
			result = true
		}
	}
	return result
}

func (h *hub) run(logger *log.Logger) {
	for {
		select {
		case c := <-h.register:
			if ct, ok := h.online[c.token]; ok {
				ct = append(ct, c)
				h.online[c.token] = ct
			} else {
				li := make([]*connection, 0)
				li = append(li, c)
				h.online[c.token] = li
			}

			mli, err := fetchOffline(c.token)
			if err != nil {
				for _, mit := range mli {
					c.send <- []byte(mit)
				}
			}

		case c := <-h.unregister:
			if _, ok := h.online[c.token]; ok {
				delete(h.online, c.token)
				close(c.send)
			}
		case m := <-h.message:
			msg, err := newMsg(m)
			if err != nil {
				logger.Printf("newMsg err:%v", err)
				break
			}
			to := msg.To
			msg.CreateTime = getFormatNow("num")
			var status bool
			status, msg = handleMsg(msg, logger)
			go func() {
				err := storeMessage(*msg)
				if err != nil {
					logger.Println("storeMessage err:%v", err)
				}

			}()

			if status {
				break
			}

			if li, ok := h.online[to]; ok {

				for _, c := range li {
					select {
					case c.send <- m:
					default:
						close(c.send)
						delete(h.online, to)
					}
				}
			}
		}
	}
}

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

func (h *hub) run(logger *log.Logger) {
	for {
		select {
		case c := <-h.register:
			if ct, ok := h.online[c.token]; ok {
				ct = append(ct, c)
			} else {
				li := make([]*connection, 0)
				li = append(li, c)
				h.online[c.token] = li
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

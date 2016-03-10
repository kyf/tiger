package main

import (
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 1024 * 1
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws *websocket.Conn

	token string

	send chan []byte
}

func (c *connection) readPump(logger *log.Logger) {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logger.Printf("error: %v", err)
			}
			break
		}
		h.message <- message
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump(logger *log.Logger) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func storeMessage(msg Message) error {
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data := make(url.Values)
	data.Set("msg", string(m))
	_, err = http.PostForm(fmt.Sprintf("%sstore", HTTP_SERVICE_URL), data)
	if err != nil {
		return err
	}
	return nil
}

func auth(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	r.ParseForm()
}

func serveWS(w http.ResponseWriter, r *http.Request, logger *log.Logger, params martini.Params) {
	token := params["token"]
	if strings.EqualFold("", token) {
		logger.Printf("token is empty")
		return
	}
	devicetoken, err := getDevicetokenByToken(token)
	if err != nil {
		logger.Printf("getDevicetokenByToken err:%v", err)
		return
	}

	if devicetoken == nil {
		logger.Printf("token %s is invalid", token)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Printf("initial websocket err:%v", err)
		return
	}
	c := &connection{token: string(devicetoken), send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump(logger)
	c.readPump(logger)
}

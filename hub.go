// Copyright 2014 The Yxorp Authors. All rights reserved.

package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Conns map[*websocket.Conn]bool

	chAdd chan *websocket.Conn
	chDel chan *websocket.Conn
	chMsg chan interface{}
}

func NewHub() *Hub {
	h := &Hub{
		Conns: make(map[*websocket.Conn]bool),
		chAdd: make(chan *websocket.Conn, 8),
		chDel: make(chan *websocket.Conn, 8),
		chMsg: make(chan interface{}, 8),
	}

	go func() {
		for {
			select {
			case c := <-h.chAdd:
				h.Conns[c] = true

			case c := <-h.chDel:
				delete(h.Conns, c)

			case v := <-h.chMsg:
				for c := range h.Conns {
					c.WriteJSON(v)
				}
			}
		}
	}()

	return h
}

func (h *Hub) Add(c *websocket.Conn) {
	h.chAdd <- c

	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(5 * time.Second))
		return nil
	})

	go func() {
		defer func() {
			h.chDel <- c
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer func() {
			c.Close()
		}()

		for {
			<-time.After(3 * time.Second)

			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}()
}

func (h *Hub) Send(v interface{}) {
	h.chMsg <- v
}

var hub = NewHub()

func handleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "", 400)
		return
	}
	catch(err)

	hub.Add(conn)

	err = conn.WriteJSON([]interface{}{"HELO"})
	catch(err)
}

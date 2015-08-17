package main

import (
	"flag"
)

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan *BroadcastMessage

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

type BroadcastMessage struct {
	Bytes []byte
	Conn  *connection
}

type connection struct {
	// The websocket connection.
	socket net.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

var h = hub{
	broadcast:   make(chan *BroadcastMessage),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			go log.Print("Register")
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			go log.Print("Unregister")
		case m := <-h.broadcast:
			for c := range h.connections {

				if reflect.DeepEqual(c, m.Conn) {
					log.Print("Skip me from broadcasting")
					continue
				}

				select {
				case c.send <- m.Bytes:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {

	defer func() {
		h.unregister <- c
		c.socket.Close()

		log.Print("Close socket on readPump")
	}()

	var buf = make([]byte, 1024)

	for {

		n, ok := c.socket.Read(buf)

		if ok != nil {
			break
		}

		log.Print("Read:", string(buf[:n]))

		message := BroadcastMessage{
			buf[:n],
			c,
		}

		h.broadcast <- &message
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {

	defer func() {
		log.Print("Close socket")
		c.socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Print("Write null byte")
				c.socket.Write([]byte{})
				return
			}
			if _, err := c.socket.Write(message); err != nil {
				log.Print("Fail to broadcast:", err)
				return
			}

			log.Print("Broadcast:", message)
		}
	}
}

func handle(s net.Conn) {

	c := &connection{
		send:   make(chan []byte, 256),
		socket: s,
	}

	h.register <- c

	go c.writePump()

	c.readPump()
}

func main() {

	port := flag.String("bind", "", "Bind address :9999 or localhost:9999")
	flag.Parse()

	if *port == "" {
		fmt.Println("Error: Please specify the bind address with --bind")
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", *port)

	if err != nil {
		panic("Fail to listen")
	}

	go h.run()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Print("Fail to accept")
			continue
		}

		go handle(conn)
	}
}

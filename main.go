package main

import (
	"bufio"
	"log"
	"net"
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
	// The socket connection.
	socket net.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

/**
 * This class handles the connections registrations/deregestration
 * broadcasting messaging and other stuff. This is mostly like an
 * event bus for wiring the communication between sockets
 */

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
			log.Print("Register")
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			log.Print("Unregister")
		case m := <-h.broadcast:
			for c := range h.connections {

				if configSocket.SkipMe && reflect.DeepEqual(c, m.Conn) {
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

/**
 * Pumps messages from the socket connection to the hub.
 */
func (c *connection) readPump() {

	defer func() {
		h.unregister <- c
		c.socket.Close()

		log.Print("Close socket on readPump")
	}()

	connbuf := bufio.NewReader(c.socket)

	for {

		str, ok := connbuf.ReadString('\n')

		if ok != nil {
			break
		}

		log.Print("Read:", str)

		message := BroadcastMessage{
			[]byte(str),
			c,
		}

		h.broadcast <- &message
	}
}

/**
 * Pumps messages from the hub to the socket connection.
 */
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

			log.Print("Broadcast:", string(message))
		}
	}
}

func HandleSocket(s net.Conn) {

	c := &connection{
		send:   make(chan []byte, 256),
		socket: s,
	}

	h.register <- c

	go c.writePump()

	c.readPump()
}

func main() {

	InitConfig()

	InitLog()

	InitMongo()

	ln, err := net.Listen("tcp", configSocket.BindAddress)

	if err != nil {
		die("Fail to listen")
	}

	socketLog.Print("Socket started on")
	apiLog.Print("Api started")

	go h.run()
	go ApiRouterInit()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Print("Fail to accept")
			continue
		}

		go HandleSocket(conn)
	}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
)

var (
	skipMe *bool
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
			go log.Print("Register")
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
			go log.Print("Unregister")
		case m := <-h.broadcast:
			for c := range h.connections {

				if *skipMe && reflect.DeepEqual(c, m.Conn) {
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

func handle(s net.Conn) {

	c := &connection{
		send:   make(chan []byte, 256),
		socket: s,
	}

	h.register <- c

	go c.writePump()

	c.readPump()
}

/**
 * Provide informations about the connections
 * by accessing the /info route on the bind port
 */
func info(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/info" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	connLen := len(h.connections)

	info := fmt.Sprintf("Connections: %s", connLen)

	w.Write([]byte(info))

	w.WriteHeader(200)
}

func main() {

	port := flag.String("bind", "", "Bind address :9999 or localhost:9999")
	skipMe = flag.Bool("skipme", true, "Skip me from broadcasting")
	portInfo := flag.String("bindinfo", "", "Bind address for informations about memory :9991 or localhost:9999")

	flag.Parse()

	if *port == "" {
		fmt.Println("Error: Please specify the bind address with --bind")
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", *port)

	if err != nil {
		panic("Fail to listen")
	}

	if *skipMe {
		log.Print("Skip me from broadcasting is ON")
	} else {
		log.Print("Skip me from broadcasting is OFF")
	}

	go h.run()

	http.HandleFunc("/info", info)
	go http.ListenAndServe(*portInfo, nil)

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Print("Fail to accept")
			continue
		}

		go handle(conn)
	}
}

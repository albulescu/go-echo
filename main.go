package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var mgop *MgoProxy

func main() {

	InitConfig()
	InitLog()

	m := &MgoProxy{}
	m.Init()

	go ApiRouterInit()
	go SocketInit()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

	fmt.Println("Graceful shutdown ...")

	m.Close()

	c := 0
	for conn, _ := range h.connections {
		conn.socket.Close()
		c++
	}

	fmt.Println("Closed socket connections:", c, "...")

	masterSocketListen.Close()
	fmt.Println("Master listen socket closed...")

	router.Shutdown()
	fmt.Println("Shutdown api router...")
}

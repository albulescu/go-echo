package main

import (
	"fmt"
)

/**
 * Provide informations about the connections
 * by accessing the /info route on the bind port
 *
 * @route /api/info
 */
func ApiInfo(request Request, response Response) {
	/*
		connLen := len(h.connections)

		info := fmt.Sprintf("Connections: %s", connLen)

		w.Write([]byte(info))

		w.WriteHeader(200)*/

	response.Error("Not implemented", 0, 500)
}

/**
 * Authentication
 *
 * @param {[type]} request  Request  [description]
 * @param {[type]} response Response [description]
 *
 * @route /api/auth
 */
func ApiAuth(request Request, response Response) {
	fmt.Print("ApiAuth")
}

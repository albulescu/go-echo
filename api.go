package main

/**
 * Provide informations about the connections
 * by accessing the /info route on the bind port
 *
 * @route /api/info
 */
func ApiInfo(request Request, response Response) {
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
	response.Error("Not implemented", 0, 500)
}

package main

import (
	"fmt"
	"log"
	"os"
)

type ErrorMessage struct {
	s string
}

func (e *ErrorMessage) Error() string {
	return e.s
}

var (
	socketLog *log.Logger
	apiLog    *log.Logger
)

func InitLog() {

	sockFile, err := os.OpenFile(configSocket.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700)

	if err != nil {
		die("Fail to open socket log file")
	} else {
		socketLog = log.New(sockFile, "", log.Ldate|log.Ltime)
	}

	apiFile, err := os.OpenFile(configApi.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0700)

	if err != nil {
		die("Fail to open api log file")
	} else {
		apiLog = log.New(apiFile, "", log.Ldate|log.Ltime)
	}
}

func die(format string, params ...interface{}) {
	fmt.Printf(format, params...)
	fmt.Print("\n")
	os.Exit(1)
}

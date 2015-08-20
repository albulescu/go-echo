package main

import (
	"flag"
	"gopkg.in/ini.v1"
	"path/filepath"
)

/**
 * This structure holds all ini configuration
 * To use new configuration from ini just add it here
 * in the structure
 */
type SocketConfiguration struct {
	// Specify bind address for web socket
	BindAddress string `ini:"bind"`
}

type ApiConfiguration struct {
	// Specify bind address for api
	BindAddress string `ini:"bind"`
}

var configSocket = new(SocketConfiguration)
var configApi = new(ApiConfiguration)

func InitConfig() {

	var configFileFlag = flag.String("config", "/etc/iapp/iapp.ini", "Config file")
	flag.Parse()

	configFile, err := filepath.Abs(*configFileFlag)

	if err != nil {
		panic("Fail to get absolute file path for ini")
	}

	ini, err := ini.Load(configFile)

	if err != nil {
		panic("Fail to load ini file")
	}

	ini.Section("socket").MapTo(configSocket)
	ini.Section("api").MapTo(configApi)
}

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

	LogFile string `ini:"log"`

	SkipMe bool `ini:"skipme"`
}

type ApiConfiguration struct {

	// Specify bind address for api
	BindAddress string `ini:"bind"`

	LogFile string `ini:"log"`

	SecretKey string `ini:"secret"`

	Expire string `ini:"expire"`
}

type DbConfiguration struct {
	Host     string `ini:"host"`
	Timeout  string `ini:"timeout"`
	Port     int    `ini:"port"`
	Username string `ini:"user"`
	Password string `ini:"pass"`
	Database string `ini:"database"`
	Source   string `ini:"source"`
}

var configSocket = new(SocketConfiguration)
var configApi = new(ApiConfiguration)
var configDb = new(DbConfiguration)

func InitConfig() {

	var configFileFlag = flag.String("config", "/etc/iapp/iapp.ini", "Config file")
	flag.Parse()

	configFile, err := filepath.Abs(*configFileFlag)

	if err != nil {
		die("Fail to get absolute file path for ini")
	}

	ini, err := ini.Load(configFile)

	if err != nil {
		die("Fail to load ini file")
	}

	ini.Section("socket").MapTo(configSocket)
	ini.Section("api").MapTo(configApi)
	ini.Section("database").MapTo(configDb)
}

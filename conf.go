package main

import (
	"fmt"
)

type config struct {
	Debug   bool
	WebHost string
	WebPort int
	DbHost string
	DbPort int
	DbName string
	Username string
	Password string
	StaticPath string
}

var Config = new(config)


func (c *config) String() string {
	s := "Config:"
	s += fmt.Sprintf("   Debug: %v\n", c.Debug)
	return s
}

func init() {
	// defaults
	Config.WebHost = "0.0.0.0"
	Config.WebPort = 8080
	Config.DbHost = "127.0.0.1"
	Config.DbPort = 0
	Config.DbName = "gotcha"
	Config.Username = "gotcha"
	Config.Password = "gotcha"
	Config.Debug = false
	Config.StaticPath = "./static"
}

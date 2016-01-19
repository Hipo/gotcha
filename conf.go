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

	StaticPath string
}

var Config = new(config)

func (c *config) HostString() string {
	return fmt.Sprintf("%s:%d", c.WebHost, c.WebPort)
}
func (c *config) DbHostString() string {
	if c.DbPort > 0 {
		return fmt.Sprintf("mongodb://%s:%d", c.DbHost, c.DbPort)
	}
	return fmt.Sprintf("mongodb://%s", c.DbHost)
}

func (c *config) String() string {
	s := "Config:"
	s += fmt.Sprintf("   Host: %s,\n", c.HostString())
	s += fmt.Sprintf("   DB: %s,\n", c.DbHostString())
	s += fmt.Sprintf("   Debug: %v\n", c.Debug)
	return s
}

func init() {
	// defaults
	Config.WebHost = "0.0.0.0"
	Config.WebPort = 8080
	Config.DbHost = "127.0.0.1"
	Config.DbPort = 0
	Config.DbName = "kriter_test"
	Config.Debug = false
	Config.StaticPath = "./static"
}

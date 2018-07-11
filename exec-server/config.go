package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Auth string `toml:"auth"`
	Bind string `toml:"bind"`
	Ca   string `toml:"ca,omitempty"`
	Cert string `toml:"cert"`
	Key  string `toml:"key"`
}

func newConfig() *Config {
	return &Config{
		Auth: "localhost:8080",
		Bind: ":8082",
		Ca:   "~/ca.pem",
		Cert: "~/service.pem",
		Key:  "~/service.key",
	}
}

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: exec-server [options]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) setFlags() *flag.FlagSet {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	fl.StringVar(&c.Auth, "auth", c.Auth, "Authentication service address.")
	fl.StringVar(&c.Bind, "bind", c.Bind, "Bind server to address.")
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")
	fl.StringVar(&c.Cert, "cert", c.Cert, "Service TLS certificate.")
	fl.StringVar(&c.Key, "key", c.Key, "Service TLS key.")

	return fl
}

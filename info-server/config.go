package main

import (
	"flag"
	"fmt"
)

type Config struct {
	User    string `toml:"user,omitempty"`
	Pass    string `toml:"pass,omitempty"`
	Token   string `toml:"token,omitempty"`
	Auth    string `toml:"auth"`
	Catalog string `toml:"catalog,omitempty"`
	Bind    string `toml:"bind"`
	Ca      string `toml:"ca,omitempty"`
	Cert    string `toml:"cert"`
	Key     string `toml:"key"`
}

func newConfig() *Config {
	return &Config{
		User:    "info-server",
		Token:   "~/.grpc-exec-example.tkn",
		Auth:    "localhost:8080",
		Catalog: "localhost:8083",
		Bind:    ":8081",
		Ca:      "../tls_setup/certs/ca.pem",
		Cert:    "../tls_setup/certs/info.pem",
		Key:     "../tls_setup/certs/info.key",
	}
}

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: info-server [options]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) setFlags() *flag.FlagSet {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	fl.StringVar(&c.User, "user", c.User, "Username for service account.")
	fl.StringVar(&c.Pass, "pass", c.Pass, "Password for service account.")
	fl.StringVar(&c.Token, "token", c.Token, "Token for service account.")
	fl.StringVar(&c.Auth, "auth", c.Auth, "Authentication service address.")
	fl.StringVar(&c.Catalog, "catalog", c.Catalog, "Catalog service address.")
	fl.StringVar(&c.Bind, "bind", c.Bind, "Bind server to address.")
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")
	fl.StringVar(&c.Cert, "cert", c.Cert, "Service TLS certificate.")
	fl.StringVar(&c.Key, "key", c.Key, "Service TLS key.")

	return fl
}

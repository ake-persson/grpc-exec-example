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
		User:    "runshit-info",
		Token:   "~/service.tkn",
		Auth:    "runshit-auth:8080",
		Catalog: "runshit-catalog:8080",
		Bind:    ":8080",
		Ca:      "~/ca.pem",
		Cert:    "~/service.pem",
		Key:     "~/service.key",
	}
}

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit-info [options]\n\nOptions:\n")
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

package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Backend    string `toml:"backend"`
	Addr       string `toml:"addr"`
	Base       string `toml:"base"`
	Domain     string `toml:"domain"`
	PrivKey    string `toml:"privKey"`
	PublKey    string `toml:"publKey"`
	Skew       int    `toml:"jwtSkew"`
	Expiration int    `toml:"expiration"`
	Bind       string `toml:"bind"`
	Ca         string `toml:"ca"`
	Cert       string `toml:"cert"`
	Key        string `toml:"key"`
}

func newConfig() *Config {
	return &Config{
		Backend:    "ldap",
		Addr:       "ldap:389",
		PrivKey:    "~/public.rsa",
		PublKey:    "~/private.rsa",
		Skew:       300,
		Expiration: 86400,
		Bind:       ":8080",
		Ca:         "~/ca.pem",
		Cert:       "~/service.pem",
		Key:        "~/service.key",
	}
}

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: auth-server [options]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) setFlags() *flag.FlagSet {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	fl.StringVar(&c.Backend, "backend", c.Backend, "Backend either ad or ldap.")
	fl.StringVar(&c.Addr, "addr", c.Addr, "LDAP server address.")
	fl.StringVar(&c.Base, "base", c.Base, "LDAP base.")
	fl.StringVar(&c.Domain, "domain", c.Domain, "LDAP domain.")
	fl.StringVar(&c.PrivKey, "priv-key", c.PrivKey, "JWT private RSA key.")
	fl.StringVar(&c.PublKey, "publ-key", c.PublKey, "JWT public RSA key.")
	fl.IntVar(&c.Skew, "skew", c.Skew, "JWT token time skew in seconds.")
	fl.IntVar(&c.Expiration, "expiration", c.Expiration, "JWT token expiration in seconds.")
	fl.StringVar(&c.Bind, "bind", c.Bind, "Bind server to address.")
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")
	fl.StringVar(&c.Cert, "cert", c.Cert, "Service TLS certificate.")
	fl.StringVar(&c.Key, "key", c.Key, "Service TLS key.")

	return fl
}

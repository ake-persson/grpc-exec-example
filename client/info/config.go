package info

import (
	"flag"
	"fmt"
)

type Config struct {
	Token   string `toml:"token,omitempty"`
	Ca      string `toml:"ca,omitempty"`
	AsJson  bool
	DefPort int
	Targets []string
}

func newConfig() *Config {
	return &Config{
		Token:  "~/service.tkn",
		Ca:     "../tls_setup/certs/ca.pem",
		AsJson: false,
	}
}

func usage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: client [options]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) setFlags() *flag.FlagSet {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = usage(fl)

	fl.StringVar(&c.Token, "token", c.Token, "Token for service account.")
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")
	fl.BoolVar(&c.AsJson, "json", c.AsJson, "Output as JSON.")
	fl.IntVar(&c.DefPort, "def-port", 8081, "Default port for info-server.")

	return fl
}

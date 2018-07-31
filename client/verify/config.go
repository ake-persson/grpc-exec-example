package verify

import (
	"flag"
	"fmt"
)

type Config struct {
	Token  string `toml:"token,omitempty"`
	Auth   string `toml:"auth"`
	Ca     string `toml:"ca,omitempty"`
	AsJson bool
}

func newConfig() *Config {
	return &Config{
		Token:  "~/service.tkn",
		Auth:   "localhost:8080",
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
	fl.StringVar(&c.Auth, "auth", c.Auth, "Authentication service address.")
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")
	fl.BoolVar(&c.AsJson, "json", c.AsJson, "Output as JSON.")

	return fl
}

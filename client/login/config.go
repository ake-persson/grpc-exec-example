package login

import (
	"flag"
	"fmt"
	"os/user"
)

type Config struct {
	User  string `toml:"user,omitempty"`
	Pass  string `toml:"pass,omitempty"`
	Token string `toml:"token,omitempty"`
	Auth  string `toml:"auth"`
	Ca    string `toml:"ca,omitempty"`
}

func newConfig() *Config {
	u, _ := user.Current()

	return &Config{
		User:  u.Username,
		Token: "~/.grpc-exec-example.tkn",
		Auth:  "localhost:8080",
		Ca:    "../tls_setup/certs/ca.pem",
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

	fl.StringVar(&c.User, "user", c.User, "Username for service account.")
	fl.StringVar(&c.Pass, "pass", c.Pass, "Password for service account.")
	fl.StringVar(&c.Token, "token", c.Token, "Token for service account.")
	fl.StringVar(&c.Auth, "auth", c.Auth, "Authentication service address.")
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")

	return fl
}

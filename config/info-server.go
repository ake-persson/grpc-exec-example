package config

import (
	"flag"
	"fmt"
	"os"
)

func infoServerUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit-info [address]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseInfoServerFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Info.User, "user", c.Info.User, "Username for service account.")
	fl.StringVar(&c.Info.Pass, "pass", c.Info.Pass, "Password for service account.")
	fl.StringVar(&c.Info.Token, "token", c.Info.Token, "Token for service account.")
	fl.StringVar(&c.Info.Catalog, "catalog", c.Info.Catalog, "Catalog address.")

	fl.StringVar(&c.Info.Bind, "bind", c.Info.Bind, "Bind server to address.")
	fl.StringVar(&c.Info.Cert, "cert", c.Info.Cert, "Server TLS certificate.")
	fl.StringVar(&c.Info.Key, "key", c.Info.Key, "Server TLS key.")

	fl.StringVar(&c.AuthClient.CA, "auth-ca", c.AuthClient.CA, "Auth. TLS CA certificate.")
	fl.BoolVar(&c.AuthClient.Insecure, "auth-insec", c.AuthClient.Insecure, "Auth. allow TLS certificates considered insecure.")
	fl.StringVar(&c.AuthClient.Address, "auth-addr", c.AuthClient.Address, "Auth. address.")

	fl.Usage = infoServerUsage(fl)

	c.parseFlags(fl, args)
}

package config

import (
	"flag"
	"fmt"
	"os"
)

func infoCatalogUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit-exec [options] [address]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseCatalogServerFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Catalog.Bind, "bind", c.Catalog.Bind, "Bind server to address.")
	fl.StringVar(&c.Catalog.Cert, "cert", c.Catalog.Cert, "Server TLS certificate.")
	fl.StringVar(&c.Catalog.Key, "key", c.Catalog.Key, "Server TLS key.")

	fl.StringVar(&c.AuthClient.CA, "auth-ca", c.AuthClient.CA, "Auth. TLS CA certificate.")
	fl.BoolVar(&c.AuthClient.Insecure, "auth-insec", c.AuthClient.Insecure, "Auth. allow TLS certificates considered insecure.")

	fl.Usage = infoCatalogUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 1 {
		if c.AuthClient.Address == "" {
			infoCatalogUsage(fl)
			os.Exit(0)
		}
	} else {
		c.AuthClient.Address = fl.Args()[0]
	}
}

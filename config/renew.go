package config

import (
	"flag"
	"fmt"
	"os"
)

func renewUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit renew [options] [address]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseRenewFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Client.Token, "token", c.Client.Token, "JWT Token.")

	fl.StringVar(&c.AuthClient.CA, "auth-ca", c.AuthClient.CA, "Auth. TLS CA certificate.")
	fl.BoolVar(&c.AuthClient.Insecure, "auth-insec", c.AuthClient.Insecure, "Auth. allow TLS certificates considered insecure.")

	fl.Usage = renewUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 1 {
		if c.AuthClient.Address == "" {
			renewUsage(fl)
			os.Exit(0)
		}
	} else {
		c.AuthClient.Address = fl.Args()[0]
	}
}

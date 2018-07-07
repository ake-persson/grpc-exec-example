package config

import (
	"flag"
	"fmt"
	"os"
)

func verifyUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit verify [options] [address]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseVerifyFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Client.Token, "token", c.Client.Token, "JWT Token.")
	fl.BoolVar(&c.Client.AsJSON, "json", c.Client.AsJSON, "Output as JSON.")

	fl.StringVar(&c.AuthClient.CA, "auth-ca", c.AuthClient.CA, "Auth. TLS CA certificate.")
	fl.BoolVar(&c.AuthClient.Insecure, "auth-insec", c.AuthClient.Insecure, "Auth. allow TLS certificates considered insecure.")

	fl.Usage = verifyUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 1 {
		if c.AuthClient.Address == "" {
			verifyUsage(fl)
			os.Exit(0)
		}
	} else {
		c.AuthClient.Address = fl.Args()[0]
	}
}

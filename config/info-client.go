package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func infoClientUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit info [options] [address,address...]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseInfoClientFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Client.CA, "ca", c.Client.CA, "TLS CA certificate.")
	fl.BoolVar(&c.Client.Insecure, "insec", c.Client.Insecure, "Allow TLS certificates considered insecure.")
	fl.StringVar(&c.Client.Token, "token", c.Client.Token, "JWT Token.")
	fl.BoolVar(&c.Client.AsJSON, "json", c.Client.AsJSON, "Output as JSON.")

	fl.Usage = infoClientUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 1 {
		if len(c.Client.Addresses) < 1 {
			infoClientUsage(fl)
			os.Exit(0)
		}
	} else {
		c.Client.Addresses = strings.Split(fl.Args()[0], ",")
	}
}

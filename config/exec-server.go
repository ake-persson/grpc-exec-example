package config

import (
	"flag"
	"fmt"
	"os"
)

func execServerUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit-exec [options] [address]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseExecServerFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Exec.Bind, "bind", c.Exec.Bind, "Bind server to address.")
	fl.StringVar(&c.Exec.Cert, "cert", c.Exec.Cert, "Server TLS certificate.")
	fl.StringVar(&c.Exec.Key, "key", c.Exec.Key, "Server TLS key.")

	fl.StringVar(&c.AuthClient.CA, "auth-ca", c.AuthClient.CA, "Auth. TLS CA certificate.")
	fl.BoolVar(&c.AuthClient.Insecure, "auth-insec", c.AuthClient.Insecure, "Auth. allow TLS certificates considered insecure.")

	fl.Usage = execServerUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 1 {
		if c.AuthClient.Address == "" {
			execServerUsage(fl)
			os.Exit(0)
		}
	} else {
		c.AuthClient.Address = fl.Args()[0]
	}
}

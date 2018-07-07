package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
)

func loginUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit login [options] [address]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseLoginFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fl.StringVar(&c.Client.Username, "username", u.Username, "Username.")
	fl.StringVar(&c.Client.Password, "password", c.Client.Password, "Password.")
	fl.StringVar(&c.Client.Token, "token", c.Client.Token, "JWT Token.")

	fl.StringVar(&c.AuthClient.CA, "auth-ca", c.AuthClient.CA, "Auth. TLS CA certificate.")
	fl.BoolVar(&c.AuthClient.Insecure, "auth-insec", c.AuthClient.Insecure, "Auth. allow TLS certificates considered insecure.")

	fl.Usage = loginUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 1 {
		if c.AuthClient.Address == "" {
			loginUsage(fl)
			os.Exit(0)
		}
	} else {
		c.AuthClient.Address = fl.Args()[0]
	}
}

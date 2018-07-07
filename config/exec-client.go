package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mickep76/runshit/strslice"
)

func execClientUsage(fl *flag.FlagSet) func() {
	return func() {
		fmt.Printf("Usage: runshit exec [options] [address,address...]\n\nOptions:\n")
		fl.PrintDefaults()
	}
}

func (c *Config) ParseExecClientFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Client.CA, "ca", c.Client.CA, "TLS CA certificate.")
	fl.BoolVar(&c.Client.Insecure, "insec", c.Client.Insecure, "Allow TLS certificates considered insecure.")
	fl.StringVar(&c.Client.Token, "token", c.Client.Token, "JWT Token.")
	fl.BoolVar(&c.Client.AsJSON, "json", c.Client.AsJSON, "Output as JSON.")

	var env strslice.StrSlice
	fl.Var(&env, "env", "Environment variable(s).")
	c.Client.Env = ([]string)(env)
	fl.StringVar(&c.Client.AsUser, "user", "", "User.")
	fl.StringVar(&c.Client.AsGroup, "group", "", "Group.")
	fl.StringVar(&c.Client.Dir, "dir", "", "Directory.")

	fl.Usage = execClientUsage(fl)

	c.parseFlags(fl, args)

	if len(fl.Args()) < 2 {
		if len(c.Client.Addresses) < 1 {
			execClientUsage(fl)
			os.Exit(0)
		}
	}

	posArgs := fl.Args()
	c.Client.Addresses = strings.Split(posArgs[0], ",")
	c.Client.Cmd = posArgs[1]
	if len(posArgs) > 2 {
		c.Client.Args = posArgs[2:]
	}
}

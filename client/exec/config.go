package exec

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Token   string `toml:"token,omitempty"`
	Ca      string `toml:"ca,omitempty"`
	AsJson  bool
	Env     StringList
	Cmd     string
	Args    []string
	AsUser  string
	AsGroup string
	InDir   string
	Addrs   []string
}

type StringList []string

func (s *StringList) String() string {
	return strings.Join([]string(*s), ",")
}

func (s *StringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func newConfig() *Config {
	return &Config{
		Token:  "~/service.tkn",
		Ca:     "~/ca.pem",
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
	fl.StringVar(&c.Ca, "ca", c.Ca, "TLS CA certificate.")
	fl.BoolVar(&c.AsJson, "json", c.AsJson, "Output as JSON.")
	fl.Var(&c.Env, "env", "Environment variable(s).")
	fl.StringVar(&c.AsUser, "user", "", "As user.")
	fl.StringVar(&c.AsGroup, "group", "", "As group.")
	fl.StringVar(&c.InDir, "dir", "", "In directory.")

	return fl
}

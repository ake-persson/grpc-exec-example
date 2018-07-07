package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mickep76/encdec"
	_ "github.com/mickep76/encdec/toml"
)

type Config struct {
	Auth       *AuthConfig       `toml:"auth"`
	LDAP       *LDAPConfig       `toml:"ldap"`
	JWT        *JWTConfig        `toml:"jwt"`
	AuthClient *AuthClientConfig `toml:"authClient"`
	Exec       *ExecConfig       `toml:"exec"`
	Info       *InfoConfig       `toml:"info"`
	Catalog    *CatalogConfig    `toml:"catalog"`
	Client     *ClientConfig     `toml:"client"`
}

type AuthConfig struct {
	Backend string `toml:"backend,omitempty"`
	Bind    string `toml:"bind,omitempty"`
	Cert    string `toml:"cert,omitempty"`
	Key     string `toml:"key,omitempty"`
}

type LDAPConfig struct {
	Address  string `toml:"address,omitempy"`
	CA       string `toml:"ca,omitempty"`
	Insecure bool   `toml:"insecure,omitempty"`
	Base     string `toml:"base,omitempty"`
	Domain   string `toml:"domain,omitempty"`
}

type JWTConfig struct {
	Expiration int    `toml:"expiration,omitempty"`
	Skew       int    `toml:"skew,omitempty"`
	PrivateKey string `toml:"privateKey,omitempty"`
	PublicKey  string `toml:"publicKey,omitempty"`
}

type AuthClientConfig struct {
	Address  string `toml:"address,omitempty"`
	CA       string `toml:"ca,omitempty"`
	Insecure bool   `toml:"insecure,omitempty"`
}

type ExecConfig struct {
	Bind string `toml:"bind,omitempty"`
	Cert string `toml:"cert,omitempty"`
	Key  string `toml:"key,omitempty"`
}

type InfoConfig struct {
	Bind string `toml:"bind,omitempty"`
	Cert string `toml:"cert,omitempty"`
	Key  string `toml:"key,omitempty"`
}

type CatalogConfig struct {
	Bind string `toml:"bind,omitempty"`
	Cert string `toml:"cert,omitempty"`
	Key  string `toml:"key,omitempty"`
}

type ClientConfig struct {
	Addresses []string `toml:"-"`
	CA        string   `toml:"ca,omitempty"`
	Insecure  bool     `toml:"insecure,omitempty"`
	Username  string   `toml:"username,omitempty"`
	Password  string   `toml:"-"`
	Token     string   `toml:"token,omitempty"`
	AsJSON    bool     `toml:"-"`
	AsUser    string   `toml:"-"`
	AsGroup   string   `toml:"-"`
	Dir       string   `toml:"-"`
	Env       []string `toml:"-"`
	Cmd       string   `toml:"-"`
	Args      []string `toml:"-"`
}

func NewConfig() *Config {
	return &Config{
		Auth: &AuthConfig{
			Backend: "ad",
			Bind:    ":8080",
			Cert:    "../tls_setup/certs/auth.pem",
			Key:     "../tls_setup/certs/auth.key",
		},
		LDAP: &LDAPConfig{
			Address:  "ldap:389",
			Insecure: false,
		},
		JWT: &JWTConfig{
			Expiration: 86400,
			Skew:       300,
			PrivateKey: "../tls_setup/certs/private.rsa",
			PublicKey:  "../tls_setup/certs/public.rsa",
		},
		AuthClient: &AuthClientConfig{
			Address:  "runshit-auth:8080",
			CA:       "../tls_setup/certs/ca.pem",
			Insecure: false,
		},
		Exec: &ExecConfig{
			Bind: ":8081",
			Cert: "../tls_setup/certs/exec.pem",
			Key:  "../tls_setup/certs/exec.key",
		},
		Info: &InfoConfig{
			Bind: ":8082",
			Cert: "../tls_setup/certs/info.pem",
			Key:  "../tls_setup/certs/info.key",
		},
		Catalog: &CatalogConfig{
			Bind: ":8083",
			Cert: "../tls_setup/certs/catalog.pem",
			Key:  "../tls_setup/certs/catalog.key",
		},
		Client: &ClientConfig{
			CA:    "../tls_setup/certs/ca.pem",
			Token: "~/.runshit.tkn",
		},
	}
}

func loadConfig(fn string, c *Config) error {
	if strings.HasPrefix(fn, "~") {
		fn = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(fn, "~"))
	}

	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		if err := encdec.FromFile("toml", fn, c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) LoadConfig() error {
	if err := loadConfig("~/.runshit.toml", c); err != nil {
		return err
	}
	return loadConfig("/etc/runshit/runshit.toml", c)
}

func (c *Config) parseFlags(fl *flag.FlagSet, args []string) {
	printConfig := fl.Bool("print-config", false, "Print config.")
	fl.Parse(args)

	if *printConfig {
		c.printConfig()
	}
}

func (c *Config) printConfig() {
	b, _ := encdec.ToBytes("toml", *c)
	fmt.Print(string(b))
	os.Exit(0)
}

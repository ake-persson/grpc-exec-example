package config

import "flag"

func (c *Config) ParseAuthFlags(args []string) {
	fl := flag.NewFlagSet("", flag.ExitOnError)

	fl.StringVar(&c.Auth.Backend, "backend", c.Auth.Backend, "Backend either ad or ldap.")
	fl.StringVar(&c.Auth.Bind, "bind", c.Auth.Bind, "Bind to address.")
	fl.StringVar(&c.Auth.Cert, "cert", c.Auth.Cert, "TLS certificate.")
	fl.StringVar(&c.Auth.Key, "key", c.Auth.Key, "TLS key")

	fl.StringVar(&c.LDAP.Address, "ldap-addr", c.LDAP.Address, "LDAP address.")
	fl.StringVar(&c.LDAP.CA, "ldap-ca", c.LDAP.CA, "TLS CA certificate.")
	fl.BoolVar(&c.LDAP.Insecure, "ldap-insec", c.LDAP.Insecure, "Allow TLS certificates considered insecure.")
	fl.StringVar(&c.LDAP.Base, "ldap-base", c.LDAP.Base, "LDAP base (ex. dc=example,dc=com.)")
	fl.StringVar(&c.LDAP.Domain, "ldap-domain", c.LDAP.Domain, "AD domain.")

	fl.IntVar(&c.JWT.Expiration, "jwt-exp", c.JWT.Expiration, "JWT expiration.")
	fl.IntVar(&c.JWT.Skew, "jwt-skew", c.JWT.Skew, "JWT skew.")
	fl.StringVar(&c.JWT.PrivateKey, "jwt-priv-key", c.JWT.PrivateKey, "JWT private key.")
	fl.StringVar(&c.JWT.PublicKey, "jwt-pub-key", c.JWT.PublicKey, "JWT public key.")

	c.parseFlags(fl, args)
}

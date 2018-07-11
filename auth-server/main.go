package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mickep76/auth"
	"github.com/mickep76/auth/jwt"
	_ "github.com/mickep76/auth/ldap"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb_auth "github.com/mickep76/grpc-exec-example/auth"
	"github.com/mickep76/grpc-exec-example/conf"
	"github.com/mickep76/grpc-exec-example/ts"
)

type server struct {
	jwt  *jwt.JWTServer
	auth auth.Conn
}

func (s *server) GetPublicKey(ctx context.Context, in *pb_auth.Empty) (*pb_auth.PublicKey, error) {
	return &pb_auth.PublicKey{Pem: s.jwt.PublicKeyPEM()}, nil
}

func (s *server) LoginUser(ctx context.Context, in *pb_auth.Login) (*pb_auth.SignedToken, error) {
	tokenUUID := uuid.New()
	log.Printf("%s request login user %s", tokenUUID, in.Username)

	u, err := s.auth.Login(in.Username, in.Password)
	if err != nil {
		err = fmt.Errorf("%s login user %s: %v", tokenUUID, in.Username, err)
		log.Print(err)
		return nil, err
	}
	u.UUID = tokenUUID

	t := s.jwt.NewToken(u)
	signed, err := s.jwt.SignToken(t)
	if err != nil {
		err = fmt.Errorf("%s sign token user %s: %v", tokenUUID, in.Username, err)
		log.Print(err)
		return nil, err
	}

	log.Printf("%s login user %s success", tokenUUID, in.Username)
	return &pb_auth.SignedToken{Token: signed}, nil
}

func (s *server) VerifyToken(ctx context.Context, in *pb_auth.SignedToken) (*pb_auth.Token, error) {
	log.Printf("verify token")

	t, err := s.jwt.ParseToken(in.Token)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	c := t.Claims.(*jwt.Claims)
	issuedAt := ts.Seconds(c.IssuedAt).Timestamp()
	expiresAt := ts.Seconds(c.ExpiresAt).Timestamp()

	log.Printf("%s verified token", c.UUID)

	return &pb_auth.Token{
		Uuid:      c.UUID,
		IssuedAt:  &issuedAt,
		ExpiresAt: &expiresAt,
		Username:  c.Username,
		Name:      c.Name,
		Mail:      c.Mail,
		Roles:     c.Roles,
		Renewed:   uint32(c.Renewed),
	}, nil
}

func (s *server) RenewToken(ctx context.Context, in *pb_auth.SignedToken) (*pb_auth.SignedToken, error) {
	log.Printf("renew token")

	t, err := s.jwt.ParseToken(in.Token)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	c := t.Claims.(*jwt.Claims)

	s.jwt.RenewToken(t)
	signed, err := s.jwt.SignToken(t)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	log.Printf("%s renewed token", c.UUID)
	return &pb_auth.SignedToken{Token: signed}, nil
}

func main() {
	c := newConfig()
	if err := conf.Load([]string{"/etc/runshit-exec.toml", "~/.runshit-exec.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, os.Args[1:], c)

	cfg := &tls.Config{
		ServerName: strings.Split(c.Addr, ":")[0], // Send SNI (Server Name Indication) for host that serves multiple aliases.
	}

	var err error
	as := &server{}
	as.auth, err = auth.Open(c.Backend, []string{c.Addr}, auth.WithTLS(cfg), auth.WithDomain(c.Domain), auth.WithBase(c.Base))
	if err != nil {
		log.Fatal(err)
	}
	defer as.auth.Close()

	as.jwt, err = jwt.NewJWTServer(jwt.RS512, time.Duration(c.Expiration)*time.Second, time.Duration(c.Skew)*time.Second, jwt.WithLoadKeys(c.PrivKey, c.PublKey))
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.Cert, c.Key)
	if err != nil {
		log.Fatalf("credentials: %v", err)
	}

	lis, err := net.Listen("tcp", c.Bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb_auth.RegisterAuthServer(s, as)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

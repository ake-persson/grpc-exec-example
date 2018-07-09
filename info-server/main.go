package main

import (
	"log"
	"net"
	"os"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb_auth "github.com/mickep76/grpc-exec-example/auth"
	"github.com/mickep76/grpc-exec-example/conf"
	pb_info "github.com/mickep76/grpc-exec-example/info"
	"github.com/mickep76/grpc-exec-example/system"
)

type server struct {
	system *pb_info.System
	jwt    *jwt.JWTClient
}

func newServer(auth string, ca string) (*server, error) {
	// Get public key from authentication service.
	k, err := pb_auth.GetPublicKey(auth, ca, false)
	if err != nil {
		return nil, err
	}

	s := &server{}

	// Create new jwt client.
	if s.jwt, err = jwt.NewJWTClient(jwt.WithPublicKey(k)); err != nil {
		return nil, err
	}

	// Cache system information.
	if s.system, err = system.Get(); err != nil {
		return nil, err
	}

	return s, nil
}

// GetSystem information.
func (s *server) GetSystem(ctx context.Context, in *pb_info.Empty) (*pb_info.System, error) {
	if err := s.jwt.AuthorizedGrpc(ctx); err != nil {
		log.Print(err)
		return nil, err
	}
	return s.system, nil
}

// Register placeholder required by interface.
func (s *server) Register(ctx context.Context, in *pb_info.System) (*pb_info.System, error) {
	return nil, nil
}

// KeepAlive placeholder required by interface.
func (s *server) KeepAlive(stream pb_info.Info_KeepAliveServer) error {
	return nil
}

// ListSystems placeholder required by interface.
func (s *server) ListSystems(ctx context.Context, in *pb_info.ListRequest) (*pb_info.SystemList, error) {
	return nil, nil
}

func main() {
	c := newConfig()
	if err := conf.Load([]string{"/etc/runshit-info.toml", "~/.runshit-info.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, os.Args, c)

	creds, err := credentials.NewServerTLSFromFile(c.Cert, c.Key)
	if err != nil {
		log.Fatalf("tls: %v", err)
	}

	lis, err := net.Listen("tcp", c.Bind)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srvr, err := newServer(c.Auth, c.Ca)
	if err != nil {
		log.Fatal("server: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb_info.RegisterInfoServer(s, srvr)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

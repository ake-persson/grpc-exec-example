package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb_auth "github.com/mickep76/runshit/auth"
	"github.com/mickep76/runshit/config"
	pb_info "github.com/mickep76/runshit/info"
	"github.com/mickep76/runshit/system"
	"github.com/mickep76/runshit/ts"
)

type server struct {
	system *pb_info.System
	jwt    *jwt.JWTClient
}

func newServer(auth string, ca string, insecure bool) (*server, error) {
	// Get public key from authentication service.
	k, err := pb_auth.GetPublicKey(auth, ca, insecure)
	if err != nil {
		return err
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
	if err := config.Load(c, []string{"/etc/runshit-info.toml", "~/.runshit-info.toml"}); err != nil {
		log.Fatal(err)
	}
	c.SetFlags()

	pubKey, err := pb_auth.GetPublicKey(c.AuthClient.Address, c.AuthClient.CA, c.AuthClient.Insecure)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile(c.Info.Cert, c.Info.Key)
	if err != nil {
		log.Fatalf("credentials: %v", err)
	}

	lis, err := net.Listen("tcp", c.Info.Bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srvr, err := NewServer(pubKey)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb_info.RegisterInfoServer(s, srvr)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	pb_auth "github.com/mickep76/grpc-exec-example/auth"
	"github.com/mickep76/grpc-exec-example/conf"
	pb_info "github.com/mickep76/grpc-exec-example/info"
	"github.com/mickep76/grpc-exec-example/system"
	"github.com/mickep76/grpc-exec-example/tlscfg"
	"github.com/mickep76/grpc-exec-example/ts"
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

func register(c *Config, cfg *tls.Config, creds credentials.PerRPCCredentials) {
	conn, err := grpc.Dial(c.Catalog,
		grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Printf("connect: %v", err)
		return
	}
	defer conn.Close()

	clnt := pb_info.NewInfoClient(conn)

	ctx := context.Background()

	s, err := system.Get()
	if err != nil {
		log.Printf("get system: %v", err)
		return
	}

	system, err := clnt.Register(ctx, s)
	if err != nil {
		log.Printf("info: %v", err)
		return
	}

	stream, err := clnt.KeepAlive(ctx)
	if err != nil {
		log.Printf("keep alive: %v", err)
		return
	}

	for {
		now := ts.Now().Timestamp()
		req := &pb_info.KeepAliveRequest{Uuid: system.Uuid, Timestamp: &now}
		log.Printf("keep alive: %v", req)
		stream.Send(req)
		time.Sleep(5 * time.Second)
	}
}

func main() {
	c := newConfig()
	if err := conf.Load([]string{"/etc/info-server.toml", "~/.info-server.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, os.Args[1:], c)

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

	if c.Register {
		tlsCfg, err := tlscfg.NewConfig(c.Ca, "", "", "", false)
		if err != nil {
			log.Fatal(err)
		}

		token, err := jwt.LoadSignedToken(c.Token)
		if err != nil {
			log.Fatal(err)
		}

		go register(c, tlsCfg, token)
	}

	if c.QRCode {
		qr, err := srvr.system.QRCode()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(qr)
	}

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

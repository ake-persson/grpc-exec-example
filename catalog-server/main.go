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

	pb_auth "github.com/mickep76/grpc-exec-example/auth"
	"github.com/mickep76/grpc-exec-example/conf"
	pb_info "github.com/mickep76/grpc-exec-example/info"
	"github.com/mickep76/grpc-exec-example/ts"
)

type server struct {
	systems map[string]*pb_info.System
	jwt     *jwt.JWTClient
}

func newServer(auth string, ca string) (*server, error) {
	// Get public key from authentication service.
	k, err := pb_auth.GetPublicKey(auth, ca, false)
	if err != nil {
		return nil, err
	}

	s := &server{
		systems: make(map[string]*pb_info.System),
	}

	// Create new jwt client.
	if s.jwt, err = jwt.NewJWTClient(jwt.WithPublicKey(k)); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *server) GetSystem(ctx context.Context, in *pb_info.Empty) (*pb_info.System, error) {
	return nil, nil
}

func (s *server) Register(ctx context.Context, in *pb_info.System) (*pb_info.System, error) {
	if err := s.jwt.AuthorizedGrpc(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	now := ts.Now().Timestamp()
	if _, ok := s.systems[in.Uuid]; ok {
		in.Updated = &now
	} else {
		in.Created = &now
	}
	in.LastSeen = &now

	s.systems[in.Uuid] = in
	return in, nil
}

func (s *server) KeepAlive(stream pb_info.Info_KeepAliveServer) error {
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("keep alive: %v", err)
			break
		}
		if err != nil {
			log.Printf("keep alive: %v", err)
			return err
		}
		log.Printf("keep alive: %v\n", resp)
		s.systems[resp.Uuid].LastSeen = resp.Timestamp
	}
	return nil
}

func (s *server) ListSystems(ctx context.Context, in *pb_info.ListRequest) (*pb_info.SystemList, error) {
	if err := s.jwt.AuthorizedGrpc(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	list := &pb_info.SystemList{}
	for _, v := range s.systems {
		list.Systems = append(list.Systems, v)
	}
	return list, nil
}

func main() {
	c := newConfig()
	if err := conf.Load([]string{"/etc/catalog-server.toml", "~/.catalog-server.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, os.Args[1:], c)

	creds, err := credentials.NewServerTLSFromFile(c.Cert, c.Key)
	if err != nil {
		log.Fatalf("credentials: %v", err)
	}

	lis, err := net.Listen("tcp", c.Bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srvr, err := newServer(c.Auth, c.Ca)
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

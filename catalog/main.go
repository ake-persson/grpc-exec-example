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

var (
	systems = make(map[string]*pb_info.System)
)

type server struct {
	system *pb_info.System
	jwt    *jwt.JWTClient
}

func NewServer(b []byte) (*server, error) {
	s := &server{}

	var err error
	if s.jwt, err = jwt.NewJWTClient(jwt.WithPublicKey(b)); err != nil {
		return nil, err
	}

	if s.system, err = system.Get(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *server) GetSystem(ctx context.Context, in *pb_info.Empty) (*pb_info.System, error) {
	if err := s.jwt.AuthorizedGrpc(ctx); err != nil {
		log.Print(err)
		return nil, err
	}
	return s.system, nil
}

func (s *server) Register(ctx context.Context, in *pb_info.System) (*pb_info.System, error) {
	if err := s.jwt.AuthorizedGrpc(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	now := ts.Now().Timestamp()
	if _, ok := systems[in.Uuid]; ok {
		in.Updated = &now
	} else {
		in.Created = &now
	}
	in.LastSeen = &now

	systems[in.Uuid] = in
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
		systems[resp.Uuid].LastSeen = resp.Timestamp
	}
	return nil
}

func (s *server) ListSystems(ctx context.Context, in *pb_info.ListRequest) (*pb_info.SystemList, error) {
	if err := s.jwt.AuthorizedGrpc(ctx); err != nil {
		log.Print(err)
		return nil, err
	}

	list := &pb_info.SystemList{}
	for _, v := range systems {
		list.Systems = append(list.Systems, v)
	}
	return list, nil
}

func main() {
	c := config.NewConfig()
	c.LoadConfig()
	c.ParseInfoServerFlags(os.Args[1:])

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

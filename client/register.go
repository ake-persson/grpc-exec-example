package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/runshit/config"
	pb_info "github.com/mickep76/runshit/info"
	"github.com/mickep76/runshit/system"
	"github.com/mickep76/runshit/tlscfg"
	"github.com/mickep76/runshit/ts"
)

func register(c *config.Config, cfg *tls.Config, creds credentials.PerRPCCredentials) {
	conn, err := grpc.Dial(c.Client.Addresses[0],
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

func registerCmd(osArgs []string) {
	c := config.NewConfig()
	c.LoadConfig()
	c.ParseInfoClientFlags(os.Args[2:])

	tlsCfg, err := tlscfg.NewConfig(c.Client.CA, "", "", "", false)
	if err != nil {
		log.Fatal(err)
	}

	token, err := jwt.LoadSignedToken(c.Client.Token)
	if err != nil {
		log.Fatal(err)
	}

	register(c, tlsCfg, token)
}

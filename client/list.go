package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/runshit/config"
	pb_info "github.com/mickep76/runshit/info"
	"github.com/mickep76/runshit/tlscfg"
)

func list(c *config.Config, cfg *tls.Config, creds credentials.PerRPCCredentials) {
	conn, err := grpc.Dial(c.Client.Addresses[0],
		grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Printf("connect: %v", err)
		return
	}
	defer conn.Close()

	clnt := pb_info.NewInfoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	list, err := clnt.ListSystems(ctx, &pb_info.ListRequest{})
	if err != nil {
		log.Printf("info: %v", err)
		return
	}

	if c.Client.AsJSON {
		b, _ := json.MarshalIndent(list.Systems, "", "  ")
		fmt.Println(string(b))
	} else {
		for _, s := range list.Systems {
			fmt.Print(s.FmtStringColor(s.Hostname))
		}
	}
}

func listCmd(osArgs []string) {
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

	list(c, tlsCfg, token)
}

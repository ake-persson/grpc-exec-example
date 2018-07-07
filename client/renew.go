package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/mickep76/runshit/auth"
	"github.com/mickep76/runshit/config"
	"github.com/mickep76/runshit/tlscfg"
)

func renewCmd(osArgs []string) {
	c := config.NewConfig()
	c.LoadConfig()
	c.ParseRenewFlags(os.Args[2:])

	tlsCfg, err := tlscfg.NewConfig(c.Client.CA, "", "", "", false)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(c.AuthClient.Address, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	if err != nil {
		log.Printf("connect: %v", err)
		return
	}
	defer conn.Close()

	clnt := pb.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if strings.HasPrefix(c.Client.Token, "~") {
		c.Client.Token = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(c.Client.Token, "~"))
	}

	b, err := ioutil.ReadFile(c.Client.Token)
	if err != nil {
		log.Fatal(err)
	}

	t, err := clnt.RenewToken(ctx, &pb.SignedToken{Token: string(b)})
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(c.Client.Token, []byte(t.Token), 0600); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote token to: %s", c.Client.Token)
}

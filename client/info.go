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

	"github.com/mickep76/runshit/color"
	"github.com/mickep76/runshit/config"
	pb_info "github.com/mickep76/runshit/info"
	"github.com/mickep76/runshit/tlscfg"
)

func info(c *config.Config, idx int, addr string, cfg *tls.Config, creds credentials.PerRPCCredentials) {
	defer wg.Done()

	conn, err := grpc.Dial(addr,
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

	s, err := clnt.GetSystem(ctx, &pb_info.Empty{})
	if err != nil {
		log.Printf("info: %v", err)
		return
	}

	if c.Client.AsJSON {
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Print(s.FmtStringColor(addr))
	}
	succ++
}

func infoCmd(osArgs []string) {
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

	for i, addr := range c.Client.Addresses {
		tot++
		wg.Add(1)
		go info(c, i, addr, tlsCfg, token)
	}

	wg.Wait()

	if len(c.Client.Addresses) > 1 {
		fmt.Fprintf(os.Stderr, "Total: %d\n%sSuccess:%s %d\n%sFailed:%s %d\n", tot, color.Green, color.Reset, succ, color.Red, color.Reset, tot-succ)
	}
}

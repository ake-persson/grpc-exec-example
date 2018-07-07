package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/mickep76/auth/jwt"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/runshit/color"
	"github.com/mickep76/runshit/config"
	pb "github.com/mickep76/runshit/exec"
	"github.com/mickep76/runshit/tlscfg"
)

func exec(c *config.Config, idx int, addr string, cfg *tls.Config, creds credentials.PerRPCCredentials, cmd string, args []string, env []string, usr string, grp string, dir string) {
	defer wg.Done()

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Printf("connect: %v", err)
		return
	}
	defer conn.Close()

	clnt := pb.NewExecCommandClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cmdUUID := uuid.New()
	stream, err := clnt.Exec(ctx, &pb.Command{Uuid: cmdUUID, Command: cmd, Arguments: args, Environment: env, User: usr, Group: grp, Directory: dir})
	if err != nil {
		log.Printf("exec: %v", err)
		return
	}

	for {
		m, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return
		}

		if c.Client.AsJSON {
			b, _ := json.MarshalIndent(m, "", "  ")
			fmt.Println(string(b))
		} else {
			if len(c.Client.Addresses) > 1 {
				fmt.Print(m.FmtStringColor(idx, addr))
			} else {
				fmt.Print(m.FmtString())
			}
		}
	}
	succ++
}

func execCmd(osArgs []string) {
	c := config.NewConfig()
	c.LoadConfig()
	c.ParseExecClientFlags(os.Args[2:])

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
		go exec(c, i, addr, tlsCfg, token, c.Client.Cmd, c.Client.Args, c.Client.Env, c.Client.AsUser, c.Client.AsGroup, c.Client.Dir)
	}

	wg.Wait()

	if len(c.Client.Addresses) > 1 {
		fmt.Fprintf(os.Stderr, "\nTotal: %d\n%sSuccess:%s %d\n%sFailed:%s %d\n", tot, color.Green, color.Reset, succ, color.Red, color.Reset, tot-succ)
	}
}

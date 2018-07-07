package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/mickep76/runshit/auth"
	"github.com/mickep76/runshit/config"
	"github.com/mickep76/runshit/tlscfg"
)

func loginCmd(osArgs []string) {
	c := config.NewConfig()
	c.LoadConfig()
	c.ParseLoginFlags(os.Args[2:])

	if c.Client.Password == "" {
		fmt.Fprintf(os.Stderr, "Password: ")
		b, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(os.Stderr, "\n")
		c.Client.Password = string(b)
	}

	tlsCfg, err := tlscfg.NewConfig(c.AuthClient.CA, "", "", "", false)
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

	t, err := clnt.LoginUser(ctx, &pb.Login{Username: c.Client.Username, Password: c.Client.Password})
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(c.Client.Token, "~") {
		c.Client.Token = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(c.Client.Token, "~"))
	}

	if err := ioutil.WriteFile(c.Client.Token, []byte(t.Token), 0600); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote token to: %s", c.Client.Token)
}

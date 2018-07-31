package info

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/grpc-exec-example/color"
	"github.com/mickep76/grpc-exec-example/conf"
	pb_info "github.com/mickep76/grpc-exec-example/info"
	"github.com/mickep76/grpc-exec-example/tlscfg"
)

var (
	wg   sync.WaitGroup
	tot  int
	succ int
)

func info(c *Config, idx int, addr string, cfg *tls.Config, creds credentials.PerRPCCredentials) {
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

	if c.AsJson {
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Print(s.FmtStringColor(addr))
	}
	succ++
}

func Cmd(args []string) {
	c := newConfig()
	if err := conf.Load([]string{"/etc/client.toml", "~/.client.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, args, c)

	if len(fl.Args()) < 1 {
		usage(fl)
		os.Exit(0)
	}
	c.Targets = strings.Split(fl.Args()[0], ",")

	tlsCfg, err := tlscfg.NewConfig(c.Ca, "", "", "", false)
	if err != nil {
		log.Fatal(err)
	}

	token, err := jwt.LoadSignedToken(c.Token)
	if err != nil {
		log.Fatal(err)
	}

	for i, addr := range c.Targets {
		tot++
		wg.Add(1)
		if strings.Contains(addr, ":") {
			go info(c, i, addr, tlsCfg, token)
		} else {
			go info(c, i, fmt.Sprintf("%s:%d", addr, c.DefPort), tlsCfg, token)
		}
	}

	wg.Wait()

	if len(c.Targets) > 1 {
		fmt.Fprintf(os.Stderr, "Total: %d\n%sSuccess:%s %d\n%sFailed:%s %d\n", tot, color.Green, color.Reset, succ, color.Red, color.Reset, tot-succ)
	}
}

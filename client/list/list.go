package list

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mickep76/auth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/grpc-exec-example/conf"
	pb_info "github.com/mickep76/grpc-exec-example/info"
	"github.com/mickep76/grpc-exec-example/tlscfg"
)

func list(c *Config, addr string, cfg *tls.Config, creds credentials.PerRPCCredentials) {
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

	systems, err := clnt.ListSystems(ctx, &pb_info.ListRequest{})
	if err != nil {
		log.Printf("list: %v", err)
		return
	}

	if c.AsJson {
		b, _ := json.MarshalIndent(systems, "", "  ")
		fmt.Println(string(b))
	} else {
		for _, s := range systems.Systems {
			fmt.Print(s.FmtStringColor(addr))
		}
	}
}

func Cmd(args []string) {
	c := newConfig()
	if err := conf.Load([]string{"/etc/client.toml", "~/.client.toml"}, c); err != nil {
		log.Fatalf("config: %v", err)
	}
	fl := c.setFlags()
	conf.ParseFlags(fl, args, c)

	tlsCfg, err := tlscfg.NewConfig(c.Ca, "", "", "", false)
	if err != nil {
		log.Fatal(err)
	}

	token, err := jwt.LoadSignedToken(c.Token)
	if err != nil {
		log.Fatal(err)
	}

	list(c, c.Catalog, tlsCfg, token)
}

package verify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb_auth "github.com/mickep76/grpc-exec-example/auth"
	"github.com/mickep76/grpc-exec-example/conf"
	"github.com/mickep76/grpc-exec-example/tlscfg"
)

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

	conn, err := grpc.Dial(c.Auth, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	if err != nil {
		log.Printf("connect: %v", err)
		return
	}
	defer conn.Close()

	clnt := pb_auth.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if strings.HasPrefix(c.Token, "~") {
		c.Token = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(c.Token, "~"))
	}

	b, err := ioutil.ReadFile(c.Token)
	if err != nil {
		log.Fatal(err)
	}

	t, err := clnt.VerifyToken(ctx, &pb_auth.SignedToken{Token: string(b)})
	if err != nil {
		log.Fatal(err)
	}

	if c.AsJson {
		b, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Print(t)
	}
}

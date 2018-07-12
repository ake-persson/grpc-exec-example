package login

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

	if c.Pass == "" {
		fmt.Fprintf(os.Stderr, "Password: ")
		b, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(os.Stderr, "\n")
		c.Pass = string(b)
	}

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

	t, err := clnt.LoginUser(ctx, &pb_auth.Login{Username: c.User, Password: c.Pass})
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(c.Token, "~") {
		c.Token = filepath.Join(os.Getenv("HOME"), strings.TrimPrefix(c.Token, "~"))
	}

	if err := ioutil.WriteFile(c.Token, []byte(t.Token), 0600); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote token to: %s", c.Token)
}

package exec

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mickep76/auth/jwt"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/mickep76/grpc-exec-example/color"
	"github.com/mickep76/grpc-exec-example/conf"
	pb_exec "github.com/mickep76/grpc-exec-example/exec"
	"github.com/mickep76/grpc-exec-example/tlscfg"
)

var (
	wg   sync.WaitGroup
	tot  int
	succ int
)

func exec(c *Config, idx int, addr string, cfg *tls.Config, creds credentials.PerRPCCredentials) {
	defer wg.Done()

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
		grpc.WithPerRPCCredentials(creds))
	if err != nil {
		log.Printf("connect: %v", err)
		return
	}
	defer conn.Close()

	clnt := pb_exec.NewExecCommandClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := clnt.Exec(ctx, &pb_exec.Command{Uuid: uuid.New(), Command: c.Cmd, Arguments: c.Args, Environment: c.Env, User: c.AsUser, Group: c.AsGroup, Directory: c.InDir})
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

		if c.AsJson {
			b, _ := json.MarshalIndent(m, "", "  ")
			fmt.Println(string(b))
		} else {
			if len(c.Targets) > 1 {
				fmt.Print(m.FmtStringColor(idx, addr))
			} else {
				fmt.Print(m.FmtString())
			}
		}
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

	if len(fl.Args()) < 2 {
		usage(fl)
		os.Exit(0)
	}

	c.Targets = strings.Split(fl.Args()[0], ",")
	c.Cmd = fl.Args()[1]
	if len(fl.Args()) > 2 {
		c.Args = fl.Args()[2:]
	}

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
		go exec(c, i, addr, tlsCfg, token)
	}

	wg.Wait()

	if len(c.Targets) > 1 {
		fmt.Fprintf(os.Stderr, "\nTotal: %d\n%sSuccess:%s %d\n%sFailed:%s %d\n", tot, color.Green, color.Reset, succ, color.Red, color.Reset, tot-succ)
	}
}
